package alarm

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/cache"
	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/cayo/pkg/model/pagination"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/toolkits/text"
	"gopkg.in/mgo.v2/bson"
)

// QueryList 列表
func (r AlarmRule) QueryList(c *gin.Context, team int) (pagination.HTTPResponse, error) {
	var (
		rules  []AlarmRule
		result pagination.HTTPResponse
	)
	qs := pagination.Pagination{
		Controller:  c,
		Table:       r.TableName(),
		SearchField: []string{"id", "name"},
		FilterField: []string{"task.id"},
		NewQs:       map[string]interface{}{"team": team},
	}.Query()
	if err := qs.Query.Limit(qs.Limit).Skip(qs.Offset).Sort("-create_time").All(&rules); err != nil {
		return result, err
	}

	if notify, err := r.notify(rules); err != nil {
		return result, err
	} else {
		rules = notify
	}

	if len(rules) == 0 {
		rules = make([]AlarmRule, 0)
	}
	result = pagination.HTTPResponse{
		Results: rules,
		Total:   qs.Total,
	}
	return result, nil
}

// Save 添加保存
func (r AlarmRule) Save() (int, error) {
	if r.ID == 0 {
		qs := bson.M{"name": r.Name, "task.id": bson.M{"$eq": r.Task.ID}}
		if err := q.Table(r.TableName()).QueryExist(qs, r.Name+"任务规则"); err != nil {
			return r.ID, err
		}
		r.ID = text.GenerateRandomToInt(6)
		r.Status = "unknown"
		r.CreateTime = time.Now()
		r.State = true
		return r.ID, global.DB.C(r.TableName()).Insert(&r)
	} else {
		return r.ID, global.DB.C(r.TableName()).Update(bson.M{"id": r.ID}, bson.M{"$set": bson.M{
			"name":      r.Name,
			"notify":    r.Notify,
			"condition": r.Condition,
			"metric":    r.Metric,
			"alarm":     r.Alarm,
		}})
	}
}

// Query 查询
func (r AlarmRule) Query() (result []int, err error) {
	var q []AlarmRule
	err = global.DB.C(r.TableName()).Find(bson.M{"task.id": bson.M{"$eq": r.Task.ID}}).All(&q)
	for _, h := range q {
		result = append(result, h.ID)
	}
	return
}

// QueryMany 查询
func (r AlarmRule) QueryMany(q interface{}) (result []AlarmRule, err error) {
	err = global.DB.C(r.TableName()).Find(q).All(&result)
	result, err = r.notify(result)
	return
}

// QueryTaskRules 查询
func (r AlarmRule) QueryTaskRules(taskID string) (result []AlarmRule, err error) {
	err = global.DB.C(r.TableName()).Find(bson.M{"task.id": bson.M{"$eq": taskID}}).All(&result)
	result, err = r.notify(result)
	return
}

// notify
func (r *AlarmRule) notify(alarmRule []AlarmRule) ([]AlarmRule, error) {
	for i := 0; i < len(alarmRule); i++ {
		var reorganise []interface{}
		for _, n := range alarmRule[i].Notify {
			var (
				contact  AlarmContact
				selector = bson.M{"id": n}
			)
			if err := q.Table(contact.TableName()).QueryOne(selector, &contact); err != nil {
				return alarmRule, err
			}
			reorganise = append(reorganise, map[string]interface{}{
				"id":       contact.ID,
				"username": contact.Username,
			})
		}
		if len(reorganise) == 0 {
			reorganise = make([]interface{}, 0)
		}
		alarmRule[i].MetaData = map[string]interface{}{
			"contact": reorganise,
		}
		for j := 0; j < len(alarmRule[i].Alarm); j++ {
			if alarmRule[i].Alarm[j].Threshold != "" {
				if GetMetricUnit(alarmRule[i].Metric) != "" {
					alarmRule[i].Alarm[j].Unit = GetMetricUnit(alarmRule[i].Metric)
				}
			}
		}
		// 查询任务告警状态
		v, err := cache.Pull(cache.GetCacheKey(alarmRule[i].Task.ID, alarmRule[i].ID))
		if err == nil && v == 1 {
			alarmRule[i].Status = "alarming"
		}
	}
	return alarmRule, nil
}

// Delete 删除
func (r AlarmRule) Delete(selector interface{}) error {
	return global.DB.C(r.TableName()).Remove(selector)
}
