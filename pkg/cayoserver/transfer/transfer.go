package transfer

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/yunling101/cayo/pkg/cache"
	"github.com/yunling101/cayo/pkg/cayoserver/judge"
	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/cayo/pkg/model/task"
	"github.com/yunling101/cayo/pkg/types"
	"github.com/yunling101/toolkits/text"
)

// transferClient
type transferClient struct {
	taskID string
	metric task.Metric
	index  []types.AlarmIndex
	task   task.Task
	nodeID int
}

// New 一个传输对象
func New(metric task.Metric, nodeID int) *transferClient {
	return &transferClient{taskID: metric.TaskId, metric: metric, nodeID: nodeID}
}

// Transfer 传输模块，
func (s *transferClient) Transfer() {
	if err := s.taskMetric(s.taskID); err != nil {
		log.Printf("query task.id %s error: %s", s.taskID, err.Error())
		return
	}
	rule, err := alarm.AlarmRule{}.QueryMany(q.M{"task.id": q.M{"$eq": s.taskID}, "state": true})
	if err != nil {
		log.Printf("query alarm rules error: %s", err.Error())
		return
	}
	for i := 0; i < len(rule); i++ {
		v := rule[i]
		s.alarmJudge(v, s.alarmValue(v.Metric))
		// 告警数据是否满足
		if v.Status == "unknown" {
			go s.repairData(v.ID)
		}
	}
}

// Saver 入库
func (s *transferClient) Saver(nid int32) {
	var metricCache task.MetricCache
	// cid := text.GenerateMd5ToString(cache.GetCacheKey(s.taskID, s.nodeID))
	err := q.Table(s.task.TableName()).QueryOne(q.M{"task_id": s.taskID}, &s.task)
	if err != nil {
		log.Printf("query task.id %s error: %s", s.taskID, err.Error())
		return
	}
	err = q.Table(metricCache.TableName()).QueryOne(q.M{"task_id": s.taskID, "status": true}, &metricCache)
	if err != nil || metricCache.ID == 0 {
		metricCache = task.MetricCache{
			ID:     text.GenerateRandomToInt(8),
			TaskId: s.taskID,
			Value:  len(s.task.ProbeNode),
			Status: true,
		}
		if err := metricCache.Save(); err != nil {
			log.Printf("save metric cache error: %s", err.Error())
			return
		}
	}

	s.metric.ID = metricCache.ID
	if err := s.metric.Saver(nid); err != nil {
		log.Printf("receipt metric insert error: %s", err.Error())
		return
	} else {
		var err error
		metricCache.Probe, err = metricCache.Reset(s.nodeID)
		if err != nil {
			log.Printf("reset metric cache error: %s", err.Error())
			return
		}
	}
	if len(s.task.ProbeNode) == len(metricCache.Probe) {
		if err := metricCache.Delete(); err != nil {
			log.Printf("delete metric cache error: %s", err.Error())
			return
		}
	}
}

// taskMetric
func (s *transferClient) taskMetric(id string) error {
	err := q.Table(s.task.TableName()).QueryOne(q.M{"task_id": id}, &s.task)
	if err != nil {
		return err
	}
	for i := 0; i < len(s.task.ProbeNode); i++ {
		if s.task.ProbeNode[i] == s.nodeID {
			s.index = append(s.index, types.AlarmIndex{
				ResponseTime: s.metric.ResponseTime, StatusCode: s.metric.StatusCode, ResponseCode: s.metric.ResponseCode,
			})
		} else {
			tableName := s.metric.GetTableName(int32(s.task.ProbeNode[i]))
			result, err := task.Metric{}.QueryOrder(tableName, q.M{"task_id": s.taskID}, 1, "-create_time")
			if err == nil && len(result) != 0 {
				v := result[len(result)-1]
				s.index = append(s.index, types.AlarmIndex{
					ResponseTime: v.ResponseTime, StatusCode: v.StatusCode, ResponseCode: v.ResponseCode,
				})
			}
		}
	}
	return nil
}

// alarmValue
func (s *transferClient) alarmValue(metric string) (value float64) {
	var dataCount, codeCount float64
	for j := 0; j < len(s.index); j++ {
		v := s.index[j]
		dataCount += v.ResponseTime
		codeCount += float64(v.StatusCode)
	}

	switch metric {
	case "ResponseTime":
		value = dataCount / float64(len(s.task.ProbeNode))
	case "ResponseCode":
		value = float64(s.index[rand.Intn(len(s.index))].ResponseCode)
	case "AvailablePercent":
		value = codeCount / float64(len(s.task.ProbeNode)) * 100
	case "AvailablePoint":
		value = codeCount / float64(len(s.task.ProbeNode))
	}
	return
}

// alarmJudge
func (s *transferClient) alarmJudge(rule alarm.AlarmRule, value float64) {
	countAlarm := 0
	k := cache.GetCacheKey(s.taskID, rule.ID)

	for _, h := range rule.Alarm {
		if h.Threshold != "" {
			threshold, err := strconv.Atoi(h.Threshold)
			if err == nil {
				switch rule.Condition {
				case ">":
					if value > float64(threshold) {
						countAlarm += s.alarmLevel(k, rule, h, value)
					}
				case "<":
					if value < float64(threshold) {
						countAlarm += s.alarmLevel(k, rule, h, value)
					}
				case "=":
					if value == float64(threshold) {
						countAlarm += s.alarmLevel(k, rule, h, value)
					}
				case ">=":
					if value >= float64(threshold) {
						countAlarm += s.alarmLevel(k, rule, h, value)
					}
				case "<=":
					if value <= float64(threshold) {
						countAlarm += s.alarmLevel(k, rule, h, value)
					}
				case "!=":
					if value != float64(threshold) {
						countAlarm += s.alarmLevel(k, rule, h, value)
					}
				}
			}
		}
	}
	// 恢复报警，此处状态存储在内存里（数据库里也有报警信息存储，后续版本可取消内存状态）
	// 报警条件等于0说明没有满足的，并查报警状态是否包含，不包含说明没有报警过，包含说明是恢复的报警（删除内存状态并修改库中报警状态）
	if countAlarm == 0 {
		if v, ok := cache.Pull(k); ok == nil && v == 1 {
			cache.Pop(k)
			t := alarm.AlarmTask{TaskID: s.taskID, RuleID: rule.ID, RecoveryTime: s.metric.CreateTime}
			if err := t.Recovery(); err != nil {
				log.Printf("task.id %s rule.id %v recovery error: %s", s.taskID, rule.ID, err.Error())
			}
		}
	}
}

// alarmLevel
func (s *transferClient) alarmLevel(k string, rule alarm.AlarmRule, level types.Alarm, current float64) int {
	// 判断连续几次，如果大于约定的值就执行 judge 模块并处理相应状态
	if s.alarmContinuity(k) >= level.Continuity {
		cond := judge.Cond{NodeID: s.nodeID, Alarm: level, Rule: rule}
		cond.Judge(current)
		// 报警后压入状态
		cache.Push(k, 1)
		// 删除后将每${Continuity}次报警一次
		cache.Delete(k)
	}
	return 1
}

// alarmContinuity
func (s *transferClient) alarmContinuity(k string) int {
	if val, err := cache.Get(k); err != nil {
		cache.Set(k, 1)
	} else {
		value := val + 1
		cache.Set(k, value)
		return value
	}
	return 1
}

// repairData 满足数据
func (s *transferClient) repairData(id int) {
	if len(s.task.ProbeNode) == len(s.index) {
		err := q.Table(alarm.AlarmRule{}.TableName()).UpdateOne(q.M{"id": id}, q.M{
			"status": "",
		})
		if err != nil {
			log.Printf("task.id %s repair data error: %s", s.task.TaskId, err.Error())
		}
	}
}
