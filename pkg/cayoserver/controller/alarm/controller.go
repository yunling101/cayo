package alarm

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/cayoserver/controller"
	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/cayo/pkg/model/settings"
	"github.com/yunling101/cayo/pkg/model/task"
	"github.com/yunling101/cayo/pkg/types"
	"github.com/yunling101/toolkits/text"
)

// AlarmController
type AlarmController struct {
	controller.BaseController
	Bind struct {
		Rule struct {
			Task   string
			Alarm  []alarm.AlarmRule
			Notify []int
		}
		State struct {
			ID    int
			State bool
		}
		Reset struct {
			ID       int
			Password string
		}
		Contact alarm.AlarmContact
		Task    task.Task
	}
}

// MetricList
func (w AlarmController) MetricList(c *gin.Context) {
	w.RenderSuccess(c, alarm.AlarmMetricList)
}

// RuleList
func (w AlarmController) RuleList(c *gin.Context) {
	query, err := alarm.AlarmRule{}.QueryList(c, 0)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, query)
}

// RuleAdd
func (w AlarmController) RuleAdd(c *gin.Context) {
	if v := w.BindJSON(c, &w.Bind.Rule); !v {
		return
	}
	if v := w.Params(c, []types.Param{
		{w.Bind.Rule.Task, "关联任务"},
		{w.Bind.Rule.Alarm, "规则描述"},
		{w.Bind.Rule.Notify, "报警联系人"},
	}); !v {
		return
	}
	err := q.Table(w.Bind.Task.TableName()).QueryOne(q.M{"task_id": w.Bind.Rule.Task}, &w.Bind.Task)
	if err != nil {
		w.RenderFail(c, "关联任务ID查询出错!")
		return
	}
	for i := 0; i < len(w.Bind.Rule.Alarm); i++ {
		w.Bind.Rule.Alarm[i].Notify = w.Bind.Rule.Notify
		w.Bind.Rule.Alarm[i].Task.ID = w.Bind.Rule.Task
		w.Bind.Rule.Alarm[i].Task.Name = w.Bind.Task.Name
		if id, err := w.Bind.Rule.Alarm[i].Save(); err != nil {
			w.RenderFail(c, err.Error())
			return
		} else {
			if !q.Contains(w.Bind.Task.Alarm, id) {
				w.Bind.Task.Alarm = append(w.Bind.Task.Alarm, id)
			}
		}
	}
	err = q.Table(w.Bind.Task.TableName()).UpdateOne(
		q.M{"id": w.Bind.Task.ID}, q.M{"alarm": w.Bind.Task.Alarm})
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

// RuleState
func (w AlarmController) RuleState(c *gin.Context) {
	if v := w.BindJSON(c, &w.Bind.State); !v {
		return
	}
	err := q.Table(alarm.AlarmRule{}.TableName()).UpdateOne(q.M{"id": w.Bind.State.ID}, q.M{
		"state": w.Bind.State.State,
	})
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

// RuleDelete
func (w AlarmController) RuleDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		w.RenderFail(c, "ID错误!")
		return
	}
	var (
		rule      alarm.AlarmRule
		alarmList []int
	)
	err = q.Table(rule.TableName()).QueryOne(q.M{"id": id}, &rule)
	if err != nil {
		w.RenderFail(c, "报警规则查询出错!")
		return
	}
	err = q.Table(w.Bind.Task.TableName()).QueryOne(q.M{"task_id": rule.Task.ID}, &w.Bind.Task)
	if err != nil {
		w.RenderFail(c, "关联任务查询出错!")
		return
	}
	for _, h := range w.Bind.Task.Alarm {
		if h != rule.ID {
			alarmList = append(alarmList, h)
		}
	}
	err = q.Table(w.Bind.Task.TableName()).UpdateOne(q.M{"task_id": rule.Task.ID}, q.M{
		"alarm": alarmList,
	})
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	err = rule.Delete(map[string]interface{}{"id": id})
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

// NotifyChannel
func (w AlarmController) NotifyChannel(c *gin.Context) {
	query := settings.Settings{}.GetUserSettings()
	result := alarm.GetNotifyChannel()
	if !query.UserLocal {
		for i := 0; i < len(result); i++ {
			if result[i]["id"] == "local" {
				result = append(result[:i], result[i+1:]...)
			}
		}
	}
	w.RenderSuccess(c, result)
}

// ContactList
func (w AlarmController) ContactList(c *gin.Context) {
	query, err := alarm.AlarmContact{}.QueryList(c)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, query)
}

// ContactAdd
func (w AlarmController) ContactAdd(c *gin.Context) {
	if err := c.ShouldBindJSON(&w.Bind.Contact); err != nil {
		w.RenderFail(c, "通知渠道不能为空!")
		return
	}
	if w.Bind.Contact.Channel == "local" {
		if v := w.Params(c, []types.Param{
			{Data: w.Bind.Contact.Username, Desc: "用户名"},
			{Data: w.Bind.Contact.Nickname, Desc: "姓名"},
			{Data: w.Bind.Contact.Password, Desc: "密码"},
			{Data: w.Bind.Contact.Role, Desc: "角色"},
			{Data: w.Bind.Contact.Email, Desc: "邮箱"},
		}); !v {
			return
		}
	} else {
		if v := w.Params(c, []types.Param{{Data: w.Bind.Contact.WebHook, Desc: "WebHook地址"}}); !v {
			return
		}
	}
	err := alarm.AlarmContact{}.Add(w.Bind.Contact)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

func (w AlarmController) ContactStat(c *gin.Context) {
	if v := w.BindJSON(c, &w.Bind.State); !v {
		return
	}
	err := q.Table(alarm.AlarmContact{}.TableName()).UpdateOne(q.M{"id": w.Bind.State.ID}, q.M{
		"is_active": w.Bind.State.State,
	})
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

func (w AlarmController) ContactReset(c *gin.Context) {
	if v := w.BindJSON(c, &w.Bind.Reset); !v {
		return
	}
	err := q.Table(alarm.AlarmContact{}.TableName()).UpdateOne(q.M{"id": w.Bind.Reset.ID}, q.M{
		"password": text.GenerateSha256ToString(w.Bind.Reset.Password),
	})
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

// ContactDelete
func (w AlarmController) ContactDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		w.RenderFail(c, "ID错误!")
		return
	}
	var rule []alarm.AlarmRule
	err = q.Table(alarm.AlarmRule{}.TableName()).QueryMany(q.M{"notify": id}, &rule)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	for i := 0; i < len(rule); i++ {
		err = q.Table(rule[i].TableName()).PullOne(q.M{"id": rule[i].ID}, q.M{"notify": id})
		if v := w.RenderIfFail(c, err); !v {
			return
		}
	}
	err = alarm.AlarmContact{}.Delete(map[string]interface{}{"id": id})
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

// NotifyContact
func (w AlarmController) NotifyContact(c *gin.Context) {
	query, err := alarm.AlarmContact{}.NotifyContact()
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, query)
}

func (w AlarmController) AlarmTask(c *gin.Context) {
	query, err := alarm.AlarmTask{}.QueryList(c)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, query)
}
