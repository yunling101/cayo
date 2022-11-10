package stat

import (
	"github.com/yunling101/cayo/pkg/cache"
	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/cayo/pkg/model/task"
	"github.com/yunling101/cayo/pkg/types"
	"github.com/yunling101/toolkits/date"
)

// Summary 汇总
func Summary() (m []q.M, err error) {
	taskTotal, err := q.Table(task.Task{}.TableName()).QueryCount(nil)
	if err != nil {
		return
	}
	var (
		alarmTask         []alarm.AlarmTask
		alarmAvailability int
		alarmResponse     int
	)
	err = q.Table(alarm.AlarmTask{}.TableName()).QueryMany(
		q.M{"status": "trigger"}, &alarmTask)
	if err != nil {
		return
	}
	for i := 0; i < len(alarmTask); i++ {
		v, err := cache.Pull(cache.GetCacheKey(alarmTask[i].TaskID, alarmTask[i].RuleID))
		if err == nil && v == 1 {
			if alarmTask[i].Metric == "可用探测点百分比" {
				alarmAvailability += 1
			} else if alarmTask[i].Metric == "响应时间" {
				alarmResponse += 1
			}
		}
	}
	m = []q.M{
		{"count": taskTotal, "name": "任务总数"},
		{"count": alarmAvailability + alarmResponse, "name": "报警总数"},
		{"count": alarmAvailability, "name": "可用率报警任务数"},
		{"count": alarmResponse, "name": "响应时间报警任务数"},
	}
	return
}

// WeekStat 周趋势
func WeekStat(day int) (m q.M, err error) {
	labels := make([]string, 0)
	data := make([]int, 0)
	alarms := make([]int, 0)
	n := date.New()

	for i := 0; i < day; i++ {
		s := n.DayStart(-(i + 1))
		c, _ := q.Table(task.Task{}.TableName()).QueryCount(
			q.M{"create_time": q.M{"$gte": s, "$lt": n.DayEnd(-(i + 1))}})
		labels = append(labels, s.Format("2006-01-02"))
		data = append(data, c)
	}
	for j := 0; j < len(data); j++ {
		s := n.DayStart(-(j + 1))
		x, _ := q.Table(alarm.AlarmTask{}.TableName()).QueryCount(
			q.M{"create_time": q.M{"$gte": s, "$lt": n.DayEnd(-(j + 1))}})
		alarms = append(alarms, x)
	}
	m = map[string]interface{}{"tasks": data, "alarm": alarms, "labels": labels}
	return
}

// WeekRatio 通道占比
func WeekRatio() (m q.M, err error) {
	notify := make(map[string]types.NotifyRatio, 0)
	for j := 0; j < len(alarm.GetLocalNotifyChannel()); j++ {
		n := alarm.GetLocalNotifyChannel()[j]
		notify[n["id"]] = types.NotifyRatio{ID: n["id"], Name: n["name"]}
	}
	notify["other"] = types.NotifyRatio{ID: "other", Name: "告警平台"}

	var alarmTask []alarm.AlarmTask
	err = q.Table(alarm.AlarmTask{}.TableName()).QueryMany(
		q.M{"status": "trigger"}, &alarmTask)
	if err != nil {
		return
	}
	var total float64
	for i := 0; i < len(alarmTask); i++ {
		for _, u := range alarmTask[i].Users {
			if v, ok := notify[u.Channel]; ok {
				notify[u.Channel] = types.NotifyRatio{
					ID: notify["other"].ID, Name: notify["other"].Name, Value: v.Value + 1,
				}
				total += 1
			} else {
				notify["other"] = types.NotifyRatio{
					ID: "other", Name: notify["other"].Name, Value: notify["other"].Value + 1,
				}
				total += 1
			}
		}
	}
	labels := make([]string, 0)
	data := make([]int, 0)
	for _, v := range notify {
		percent := (v.Value / total) * 100
		labels = append(labels, v.Name)
		if total == 0 {
			data = append(data, 0)
		} else {
			data = append(data, int(percent))
		}
	}

	m = map[string]interface{}{"data": data, "labels": labels}
	return
}
