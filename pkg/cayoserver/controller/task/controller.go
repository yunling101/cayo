package task

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/cayoserver/controller"
	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/cayo/pkg/model/task"
	"github.com/yunling101/cayo/pkg/types"
)

// TaskController
type TaskController struct {
	controller.BaseController
	Bind struct {
		Task struct {
			Data  task.Task
			Alarm []alarm.AlarmRule
		}
		State struct {
			ID    int
			State bool
		}
	}
}

// List
func (w TaskController) List(c *gin.Context) {
	query, err := task.Task{}.QueryList(c)
	if err != nil {
		w.RenderFail(c, err.Error())
		return
	}
	w.RenderSuccess(c, query)
}

// TaskType
func (w TaskController) TaskType(c *gin.Context) {
	w.RenderSuccess(c, task.GetTaskTypeList())
}

// DnsType
func (w TaskController) DnsType(c *gin.Context) {
	w.RenderSuccess(c, task.GetDnsTypeList())
}

// Add
func (w TaskController) Add(c *gin.Context) {
	if err := c.ShouldBindJSON(&w.Bind.Task); err != nil {
		w.RenderDieError(c, err, "参数解析错误!")
		return
	}
	if v := w.Params(c, []types.Param{
		{Data: w.Bind.Task.Data.Name, Desc: "任务名称"},
		{Data: w.Bind.Task.Data.Protocol, Desc: "任务类型"},
		{Data: w.Bind.Task.Data.Address, Desc: "监控域名"},
	}); !v {
		return
	}
	switch w.Bind.Task.Data.Protocol {
	case "TCP", "UDP":
		if v := w.Params(c, []types.Param{{Data: w.Bind.Task.Data.Port, Desc: "端口"}}); !v {
			return
		}
	case "SMTP":
		if v := w.Params(c, []types.Param{
			{Data: w.Bind.Task.Data.Port, Desc: "端口"},
			{Data: w.Bind.Task.Data.Username, Desc: "用户名"},
			{Data: w.Bind.Task.Data.Password, Desc: "密码"},
		}); !v {
			return
		}
	}
	if v := w.Params(c, []types.Param{
		{Data: w.Bind.Task.Data.ProbeNode, Desc: "探测节点"},
		{Data: w.Bind.Task.Alarm, Desc: "规则描述"},
		{Data: w.Bind.Task.Data.Notify, Desc: "通知联系人"},
	}); !v {
		return
	}
	err := task.Task{}.Add(w.Bind.Task.Data, w.Bind.Task.Alarm)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

// State
func (w TaskController) State(c *gin.Context) {
	if v := w.BindJSON(c, &w.Bind.State); !v {
		return
	}
	err := q.Table(task.Task{}.TableName()).UpdateOne(q.M{"id": w.Bind.State.ID}, q.M{
		"status": w.Bind.State.State,
	})
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

// Delete
func (w TaskController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		w.RenderFail(c, "ID错误!")
		return
	}
	err = task.Task{ID: id}.Delete()
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, id)
}
