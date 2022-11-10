package task

import (
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/toolkits/date"
)

func (t MetricCache) Save() error {
	t.CreateTime = date.New().Now
	return q.Table(t.TableName()).InsertOne(&t)
}

func (t MetricCache) Reset(nid int) (probe []int, err error) {
	for _, i := range t.Probe {
		if i == nid {
			return
		}
	}
	t.Probe = append(t.Probe, nid)
	return t.Probe, q.Table(t.TableName()).UpdateOne(q.M{"task_id": t.TaskId}, q.M{
		"probe": t.Probe,
	})
}

func (t MetricCache) Update(selector interface{}) error {
	return q.Table(t.TableName()).UpdateOne(q.M{"task_id": t.TaskId}, selector)
}

func (t MetricCache) Delete() error {
	return q.Table(t.TableName()).DeleteOne(q.M{"task_id": t.TaskId})
}
