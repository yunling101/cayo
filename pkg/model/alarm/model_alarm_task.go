package alarm

import (
	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/cache"
	"github.com/yunling101/cayo/pkg/model/pagination"
	"github.com/yunling101/cayo/pkg/model/q"
)

// QueryList 列表
func (t AlarmTask) QueryList(c *gin.Context) (pagination.HTTPResponse, error) {
	var (
		alarmTask  []AlarmTask
		alarmTasks []AlarmTask
		result     pagination.HTTPResponse
	)
	qs := pagination.Pagination{
		Controller:  c,
		Table:       t.TableName(),
		FilterField: []string{"task_name", "task_id", "rule_id", "status"},
		NumberField: []string{"rule_id"},
	}.Query()
	if err := qs.Query.Limit(qs.Limit).Skip(qs.Offset).Sort("-create_time").All(&alarmTask); err != nil {
		return result, err
	}
	for i := 0; i < len(alarmTask); i++ {
		v, err := cache.Pull(cache.GetCacheKey(alarmTask[i].TaskID, alarmTask[i].RuleID))
		if err == nil && v == 1 {
			alarmTasks = append(alarmTasks, alarmTask[i])
		}
	}
	if len(alarmTasks) == 0 {
		alarmTasks = make([]AlarmTask, 0)
		qs.Total = 0
	}
	result = pagination.HTTPResponse{
		Results: alarmTasks,
		Total:   qs.Total,
	}
	return result, nil
}

func (t AlarmTask) Save() (err error) {
	var alarm AlarmTask
	err = q.Table(t.TableName()).QueryOne(q.M{"task_id": t.TaskID, "rule_id": t.RuleID}, &alarm)
	if err != nil {
		err = q.Table(t.TableName()).InsertOne(&t)
	} else {
		if alarm.Status == "recovery" {
			err = q.Table(t.TableName()).UpdateOne(q.M{"id": alarm.ID}, q.M{
				"archive":      true,
				"archive_time": t.CreateTime,
			})
			err = q.Table(t.TableName()).InsertOne(&t)
		} else {
			err = q.Table(t.TableName()).UpdateOne(q.M{"id": alarm.ID}, q.M{
				"status":  "trigger",
				"counter": alarm.Counter + 1,
			})
		}
	}
	return
}

func (t AlarmTask) Recovery() error {
	return q.Table(t.TableName()).UpdateOne(q.M{"task_id": t.TaskID, "rule_id": t.RuleID}, q.M{
		"status":        "recovery",
		"recovery_time": t.RecoveryTime,
	})
}
