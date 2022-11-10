package client

import (
	"fmt"
	"log"

	//"strconv"
	"encoding/base64"

	jsoniter "github.com/json-iterator/go"
	"github.com/yunling101/cayo/pkg/cayoserver/transfer"
	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/cayo/pkg/model/node"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/cayo/pkg/model/task"
	"github.com/yunling101/cayo/pkg/propb"
	"github.com/yunling101/toolkits/date"
)

type client struct {
	request *propb.Request
	jsonNew jsoniter.API
	task    task.Task
}

func New(r *propb.Request) *client {
	return &client{
		request: r, jsonNew: jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

// Heartbeat 心跳监测
func (c *client) Heartbeat() (err error) {
	var endpoint node.EndPoint
	err = q.Table(endpoint.TableName()).QueryOne(q.M{"id": c.request.Nid}, &endpoint)
	if err != nil {
		log.Printf("nid %v query error: %s", c.request.Nid, err.Error())
		return
	}
	if endpoint.Status == "Pending" {
		err = q.Table(endpoint.TableName()).UpdateOne(q.M{"id": c.request.Nid},
			q.M{
				"version":   c.request.Param["version"],
				"status":    "Success",
				"heartbeat": date.New().Now,
			},
		)
	} else {
		err = q.Table(endpoint.TableName()).UpdateOne(q.M{"id": c.request.Nid},
			q.M{
				"version":   c.request.Param["version"],
				"heartbeat": date.New().Now,
			},
		)
	}
	if err != nil {
		err = fmt.Errorf("nid %d heartbeat fail: %s", c.request.Nid, err.Error())
	}
	return
}

// ObtainTask 任务过去
func (c *client) ObtainTask() (result map[string]string, err error) {
	var (
		tasks  []task.Task
		offset int = 0
	)
	limit := global.Config().NodeMaxLimit
	cond := q.M{"probe_node": c.request.Nid, "status": true}
	if err = q.Table(c.task.TableName()).QueryPage(limit, offset, cond, &tasks); err != nil {
		err = fmt.Errorf("nid %d obtain task fail: %s", c.request.Nid, err.Error())
		return
	}
	b, _ := c.jsonNew.Marshal(tasks)
	result = map[string]string{
		"data": base64.StdEncoding.EncodeToString(b),
	}
	return
}

// ReceiptMetric 回执处理
func (c *client) ReceiptMetric() {
	var (
		metric []interface{}
		point  node.EndPoint
	)
	if err := c.jsonNew.Unmarshal([]byte(c.request.Param["data"]), &metric); err != nil {
		log.Printf("receipt metric unmarshal data error: %s", err.Error())
		return
	}
	err := q.Table(point.TableName()).QueryOne(q.M{"id": int(c.request.Nid)}, &point)
	if err != nil {
		log.Printf("receipt metric query node error: %s", err.Error())
		return
	}

	// 遍历接收到的数据
	for i := 0; i < len(metric); i++ {
		// 解析返回数据
		if m, err := c.parse(metric[i]); err != nil {
			log.Printf("receipt metric parse error: %s", err.Error())
			continue
		} else {
			m.ProbePoint = fmt.Sprintf("%s-%s-%s", point.Province, point.City, node.GetOperatorName(point.Operator))
			m.ProbeSource = point.IP
			m.TaskId = c.request.Param["task_id"]
			m.CreateTime = date.New().Now
			trans := transfer.New(m, point.ID)

			// 分发数据到judge模块（异步）
			go trans.Transfer()
			// 根据任务所在节点存储历史数据
			go trans.Saver(c.request.Nid)
		}
	}
}

// parse
func (c *client) parse(v interface{}) (metric task.Metric, err error) {
	b, _ := c.jsonNew.Marshal(v)
	err = c.jsonNew.Unmarshal(b, &metric)
	return
}
