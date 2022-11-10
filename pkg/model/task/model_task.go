package task

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/cache"
	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/model/pagination"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/cayo/pkg/types"
	"github.com/yunling101/toolkits/text"
	"gopkg.in/mgo.v2/bson"
)

// QueryList 列表
func (t Task) QueryList(c *gin.Context) (pagination.HTTPResponse, error) {
	var (
		task   []Task
		result pagination.HTTPResponse
	)
	qs := pagination.Pagination{
		Controller:  c,
		Table:       t.TableName(),
		SearchField: []string{"name", "task_id"},
		FilterField: []string{"protocol"},
	}.Query()
	if err := qs.Query.Limit(qs.Limit).Skip(qs.Offset).Sort("-create_time").All(&task); err != nil {
		return result, err
	}
	result = pagination.HTTPResponse{
		Results: task,
		Total:   qs.Total,
	}
	if len(task) == 0 {
		result.Results = make([]string, 0)
	} else {
		// 查询任务告警状态
		for i := 0; i < len(task); i++ {
			for j := 0; j < len(task[i].Alarm); j++ {
				v, err := cache.Pull(cache.GetCacheKey(task[i].TaskId, task[i].Alarm[j]))
				if err == nil && v == 1 {
					task[i].AlarmState = true
				}
			}
		}
	}
	return result, nil
}

// Add 添加
func (t Task) Add(n Task, rule []alarm.AlarmRule) error {
	action := n.ID
	if action == 0 {
		if err := q.Table(t.TableName()).QueryExist(
			bson.M{"name": n.Name}, "任务名称"); err != nil {
			return err
		}
		n.ID = text.GenerateRandomToInt(6)
		n.TaskId = text.GenerateRandomUUID()
	}
	for j := 0; j < len(rule); j++ {
		rule[j].Task = types.Model{ID: n.TaskId, Name: n.Name}
		rule[j].Notify = n.Notify
	}
	var currentAlarm []Factory
	for _, h := range rule {
		currentAlarm = append(currentAlarm, h)
	}
	err := NewModal(&alarm.AlarmRule{}).Add(currentAlarm, &n.Alarm)
	if err != nil {
		return err
	}

	if action == 0 {
		n.CreateTime = time.Now()
		n.Status = true
		return global.DB.C(t.TableName()).Insert(&n)
	} else {
		return global.DB.C(t.TableName()).Update(bson.M{"id": action}, bson.M{"$set": map[string]interface{}{
			"name":            n.Name,
			"address":         n.Address,
			"pack_number":     n.PackNumber,
			"frequency":       n.Frequency,
			"port":            n.Port,
			"retry":           n.Retry,
			"method":          n.Method,
			"resolution_type": n.ResolutionType,
			"server":          n.Server,
			"tls":             n.Tls,
			"header":          n.Header,
			"username":        n.Username,
			"password":        n.Password,
			"ssl_verify":      n.SSLVerify,
			"probe_node":      n.ProbeNode,
			"alarm":           n.Alarm,
			"notify":          n.Notify,
		}})
	}
}

// Delete 删除
func (t Task) Delete() (err error) {
	var task Task
	if err = q.Table(t.TableName()).QueryOne(q.M{"id": t.ID}, &task); err != nil {
		return
	}
	for _, h := range task.Alarm {
		err = alarm.AlarmRule{}.Delete(q.M{"id": h})
		if err != nil {
			return err
		}
	}
	err = q.Table(t.TableName()).DeleteOne(q.M{"id": t.ID})
	return
}
