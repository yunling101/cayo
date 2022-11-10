package alarm

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/model/pagination"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/cayo/pkg/types"
	"github.com/yunling101/toolkits/text"
)

// QueryList 列表
func (t AlarmContact) QueryList(c *gin.Context) (pagination.HTTPResponse, error) {
	var (
		contact []AlarmContact
		result  pagination.HTTPResponse
	)
	qs := pagination.Pagination{
		Controller:  c,
		Table:       t.TableName(),
		SearchField: []string{"name", "phone", "email", "webhook"},
	}.Query()
	if err := qs.Query.Limit(qs.Limit).Skip(qs.Offset).Sort("-create_time").All(&contact); err != nil {
		return result, err
	}
	result = pagination.HTTPResponse{
		Results: contact,
		Total:   qs.Total,
	}
	if len(contact) == 0 {
		result.Results = make([]string, 0)
	}
	return result, nil
}

// Add
func (t AlarmContact) Add(n AlarmContact) error {
	if n.ID == 0 {
		if n.Channel == "local" {
			if err := q.Table(t.TableName()).QueryExist(q.M{"username": n.Username}, "用户名"); err != nil {
				return err
			}
		} else {
			if n.WebHook != "" {
				if err := q.Table(t.TableName()).QueryExist(q.M{"webhook": n.WebHook}, "WebHook"); err != nil {
					return err
				}
			}
		}
		n.ID = text.GenerateRandomToInt(7)
		n.Password = text.GenerateSha256ToString(n.Password)
		n.IsActive = true
		if n.Role == "admin" {
			n.IsSuperuser = true
		}
		n.CreateTime = time.Now()
		return q.Table(t.TableName()).InsertOne(&n)
	} else {
		return q.Table(t.TableName()).UpdateOne(q.M{"id": n.ID}, q.M{
			"username": n.Username,
			"nickname": n.Nickname,
			"channel":  n.Channel,
			"role":     n.Role,
			"phone":    n.Phone,
			"email":    n.Email,
			"dingding": n.DingDing,
			"feishu":   n.Feishu,
			"webhook":  n.WebHook,
		})
	}
}

// Delete
func (t AlarmContact) Delete(selector interface{}) error {
	return q.Table(t.TableName()).DeleteOne(selector)
}

// NotifyContact
func (t AlarmContact) NotifyContact() ([]types.Menu, error) {
	var (
		menu    []types.Menu
		contact []AlarmContact
	)
	if err := q.Table(t.TableName()).QueryMany(nil, &contact); err != nil {
		return menu, err
	}
	for _, c := range contact {
		menu = append(menu, types.Menu{
			Label: fmt.Sprintf("%s", c.Username),
			Value: c.ID,
		})
	}
	return menu, nil
}
