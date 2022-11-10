package task

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/cayo/pkg/model/pagination"
	"github.com/yunling101/toolkits/date"
	"github.com/yunling101/toolkits/text"
	"gopkg.in/mgo.v2/bson"
)

// GetTableName 获取表名称
func (t *Metric) GetTableName(nid int32) string {
	return fmt.Sprintf(t.TableName(), fmt.Sprintf("%v", nid))
}

// QueryList 列表
func (t Metric) QueryList(c *gin.Context, nid int32, s *date.Date) (metric []Metric, total int, err error) {
	qs := pagination.Pagination{
		Controller:  c,
		Table:       t.GetTableName(nid),
		SearchField: []string{"name", "task_id"},
		FilterField: []string{"task_id", "status_code"},
		NewQs: map[string]interface{}{
			"create_time": bson.M{"$gte": s.BeforeToHour(6), "$lt": s.Now},
		},
	}.Query()
	if errs := qs.Query.Limit(qs.Limit).Skip(qs.Offset).Sort("-create_time").All(&metric); errs != nil {
		err = errs
		return
	}
	total = qs.Total
	return
}

func (t Metric) Saver(nid int32) error {
	count, err := global.DB.C(t.GetTableName(nid)).Find(bson.M{"id": t.ID}).Count()
	if err != nil || count != 0 {
		t.ID = text.GenerateRandomToInt(8)
	}
	return global.DB.C(t.GetTableName(nid)).Insert(&t)
}

// QueryOrder 排序查询
func (t Metric) QueryOrder(tableName string, selector interface{}, n int, order string) (r []Metric, err error) {
	err = global.DB.C(tableName).Find(selector).Limit(n).Sort(order).All(&r)
	return
}
