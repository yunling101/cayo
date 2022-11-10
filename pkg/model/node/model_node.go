package node

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/cayo/pkg/model/pagination"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/cayo/pkg/model/task"
	"github.com/yunling101/cayo/pkg/types"
	"github.com/yunling101/toolkits/text"
)

// QueryList 列表
func (t EndPoint) QueryList(c *gin.Context) (pagination.HTTPResponse, error) {
	var (
		hosts  []EndPoint
		result pagination.HTTPResponse
	)
	qs := pagination.Pagination{
		Controller:  c,
		Table:       t.TableName(),
		SearchField: []string{"operator", "province", "city", "ip"},
	}.Query()
	if err := qs.Query.Limit(qs.Limit).Skip(qs.Offset).Sort("-create_time").All(&hosts); err != nil {
		return result, err
	}
	for i := 0; i < len(hosts); i++ {
		hosts[i].OperatorZh = GetOperatorName(hosts[i].Operator)
		c, _ := q.Table(task.Task{}.TableName()).QueryCount(q.M{"probe_node": hosts[i].ID, "status": true})
		hosts[i].TasksCount = c
	}
	result = pagination.HTTPResponse{
		Results: hosts,
		Total:   qs.Total,
	}
	if len(hosts) == 0 {
		result.Results = make([]string, 0)
	}
	return result, nil
}

// Add 添加
func (t EndPoint) Add(n EndPoint) error {
	if n.ID == 0 {
		if err := q.Table(t.TableName()).QueryExist(q.M{"ip": n.IP}, "IP地址"); err != nil {
			return err
		}
		n.ID = text.GenerateRandomToInt(6)
		n.Attribute = "私有"
		n.Status = "Pending"
		n.Selected = false
		n.Heartbeat = time.Now()
		n.CreateTime = time.Now()
		return global.DB.C(t.TableName()).Insert(&n)
	}
	return nil
}

// UpdateOne 更新
func (t EndPoint) UpdateOne(selector interface{}, update interface{}) error {
	return global.DB.C(t.TableName()).Update(selector, q.M{"$set": update})
}

// Delete 删除
func (t EndPoint) Delete(selector interface{}) error {
	return global.DB.C(t.TableName()).Remove(selector)
}

// Probe 任务节点
func (t EndPoint) Probe() ([]types.Menu, []int, error) {
	var node []EndPoint

	defaultMenu := make([]types.Menu, 0)
	selectedMenu := make([]int, 0)
	if err := global.DB.C(t.TableName()).Find(q.M{"status": "Success"}).All(&node); err != nil {
		return defaultMenu, selectedMenu, err
	}
	for _, h := range node {
		// 节点最大任务数限制
		if t.nodeTaskCount(h.ID) {
			defaultMenu = append(defaultMenu, types.Menu{
				Label: fmt.Sprintf("%s-%s-%s", GetOperatorName(h.Operator), h.Province, h.City),
				Value: h.ID,
			})
			if h.Selected {
				selectedMenu = append(selectedMenu, h.ID)
			}
		}
	}
	return defaultMenu, selectedMenu, nil
}

// nodeTaskCount 判断节点任务数
func (t EndPoint) nodeTaskCount(id int) bool {
	c, err := global.DB.C(task.Task{}.TableName()).Find(q.M{"probe_node": id}).Count() // "status": true
	if err != nil {
		return true
	}
	if c <= global.Config().NodeMaxLimit {
		return true
	}
	return false
}
