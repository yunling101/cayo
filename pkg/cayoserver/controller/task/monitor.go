package task

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/cayoserver/chart"
	"github.com/yunling101/cayo/pkg/cayoserver/controller"
	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/cayo/pkg/model/task"
	"github.com/yunling101/cayo/pkg/types"
	"github.com/yunling101/toolkits/date"
)

// MonitorController 控制器
type MonitorController struct {
	controller.BaseController
	Bind struct {
		task   task.Task
		metric task.Metric
	}
}

type metricSummary struct {
	TableName string
	Data      []task.Metric
}

// SummaryArea 分析统计
func (w MonitorController) SummaryArea(c *gin.Context) {
	uuid := c.Param("id")
	if v := w.Params(c, []types.Param{{uuid, "id"}}); !v {
		return
	}

	interval, err := strconv.Atoi(c.Query("interval"))
	if err != nil {
		interval = 1
	}
	err = q.Table(w.Bind.task.TableName()).QueryOne(q.M{"task_id": uuid}, &w.Bind.task)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	rules, err := alarm.AlarmRule{}.QueryTaskRules(uuid)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	n := chart.NewVariable()

	maxMetric := 0
	metrics := make(map[int]metricSummary, 0)
	// 因每个Agent上报的时间和条数不一致，所以需要取出个节点的数据量
	for i := 0; i < len(w.Bind.task.ProbeNode); i++ {
		tableName := w.Bind.metric.GetTableName(int32(w.Bind.task.ProbeNode[i]))
		cond := q.M{"task_id": uuid, "create_time": q.M{"$gte": n.Time.BeforeToHour(interval), "$lt": n.Time.Now}}

		var result []task.Metric
		if err = q.Table(tableName).QueryMany(cond, &result); err != nil {
			continue
		}

		if i == 0 {
			maxMetric = len(result)
		} else {
			if maxMetric < len(result) {
				maxMetric = len(result)
			}
		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].CreateTime.Unix() < result[j].CreateTime.Unix()
		})
		metrics[w.Bind.task.ProbeNode[i]] = metricSummary{TableName: tableName, Data: result}
	}

	firstMetric := false
	for _, first := range metrics {
		if maxMetric == len(first.Data) {
			if !firstMetric {
				for x := 0; x < maxMetric; x++ {
					v := first.Data[x]
					k := fmt.Sprintf("%v", first.Data[x].ID)
					n.Labels = append(n.Labels, v.CreateTime.Local().Format("2006-01-02 15:04:00"))
					n.Keys = append(n.Keys, k)

					n.AddAvail(k, v.StatusCode)
					if w.Bind.task.Protocol == "DNS" || w.Bind.task.Protocol == "PING" {
						n.AddData(k, chart.Sets{Value: map[string]float64{
							"响应时间": v.ResponseTime,
						}},
						)
					} else {
						n.AddData(k, chart.Sets{Value: map[string]float64{
							"响应时间":  v.ResponseTime,
							"DNS时间": v.DNSTime,
							"建连时间":  v.ConnectTime,
							"首包时间":  v.PackTime,
							"SSL时间": v.SSLTime,
						}},
						)
					}
					n.Results = append(n.Results, first.Data[x])
				}
				firstMetric = true
			} else {
				for x := 0; x < maxMetric; x++ {
					v := first.Data[x]
					k := fmt.Sprintf("%v", first.Data[x].ID)

					if n.GetAvail(k) {
						n.AddAvail(k, v.StatusCode)
						if w.Bind.task.Protocol == "DNS" || w.Bind.task.Protocol == "PING" {
							n.AddData(k, chart.Sets{Value: map[string]float64{
								"响应时间": v.ResponseTime,
							}},
							)
						} else {
							n.AddData(k, chart.Sets{Value: map[string]float64{
								"响应时间":  v.ResponseTime,
								"DNS时间": v.DNSTime,
								"建连时间":  v.ConnectTime,
								"首包时间":  v.PackTime,
								"SSL时间": v.SSLTime,
							}},
							)
						}
					}
					n.Results = append(n.Results, first.Data[x])
				}
			}
		}
	}

	for _, b := range metrics {
		if maxMetric != len(b.Data) {
			for x := 0; x < maxMetric; x++ {
				if x < len(b.Data) {
					v := b.Data[x]
					k := fmt.Sprintf("%v", b.Data[x].ID)

					if n.GetAvail(k) {
						n.AddAvail(k, v.StatusCode)
						if w.Bind.task.Protocol == "DNS" || w.Bind.task.Protocol == "PING" {
							n.AddData(k, chart.Sets{Value: map[string]float64{
								"响应时间": v.ResponseTime,
							}},
							)
						} else {
							n.AddData(k, chart.Sets{Value: map[string]float64{
								"响应时间":  v.ResponseTime,
								"DNS时间": v.DNSTime,
								"建连时间":  v.ConnectTime,
								"首包时间":  v.PackTime,
								"SSL时间": v.SSLTime,
							}},
							)
						}
					}
					n.Results = append(n.Results, b.Data[x])
				} else {
					k := n.Keys[x]

					n.AddAvail(k, 0)
					if w.Bind.task.Protocol == "DNS" || w.Bind.task.Protocol == "PING" {
						n.AddData(k, chart.Sets{Value: map[string]float64{
							"响应时间": 0,
						}},
						)
					} else {
						n.AddData(k, chart.Sets{Value: map[string]float64{
							"响应时间":  0,
							"DNS时间": 0,
							"建连时间":  0,
							"首包时间":  0,
							"SSL时间": 0,
						}},
						)
					}
				}
			}
		}
	}
	n.Node = len(w.Bind.task.ProbeNode)
	w.RenderSuccess(c, chart.Get(n, rules))
}

// SummaryMetric 探测结果
func (w MonitorController) SummaryMetric(c *gin.Context) {
	uuid := c.Query("task_id")
	if v := w.Params(c, []types.Param{{uuid, "id"}}); !v {
		return
	}

	err := q.Table(w.Bind.task.TableName()).QueryOne(q.M{"task_id": uuid}, &w.Bind.task)
	if err != nil {
		w.RenderFail(c, err.Error())
		return
	}

	var (
		results []task.Metric
		total   int
	)
	s := date.New()
	for i := 0; i < len(w.Bind.task.ProbeNode); i++ {
		result, num, err := task.Metric{}.QueryList(c, int32(w.Bind.task.ProbeNode[i]), s)
		if err == nil {
			if total < num {
				total = num
			}
			results = append(results, result...)
		}
	}
	if len(results) == 0 {
		results = make([]task.Metric, 0)
	} else {
		sort.Slice(results, func(i, j int) bool {
			return results[i].CreateTime.Unix() > results[j].CreateTime.Unix()
		})
	}
	w.RenderSuccess(c, map[string]interface{}{"results": results, "total": total})
}

// SummaryRules 报警规则
func (w MonitorController) SummaryRules(c *gin.Context) {
	uuid := c.Query("task_id")
	if v := w.Params(c, []types.Param{{uuid, "id"}}); !v {
		return
	}
	rules, err := alarm.AlarmRule{}.QueryTaskRules(uuid)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	if len(rules) == 0 {
		rules = make([]alarm.AlarmRule, 0)
	}
	w.RenderSuccess(c, rules)
}
