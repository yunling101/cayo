package console

import (
	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/auth"
	"github.com/yunling101/cayo/pkg/cayoserver/controller"
	"github.com/yunling101/cayo/pkg/model/stat"
)

// ConsoleController
type ConsoleController struct {
	controller.BaseController
}

// Owner 本人信息
func (w ConsoleController) Owner(c *gin.Context) {
	w.RenderSuccess(c, auth.Request.User)
}

// Summary 汇总
func (w ConsoleController) Summary(c *gin.Context) {
	data, err := stat.Summary()
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, data)
}

// WeekStat 近期数据
func (w ConsoleController) WeekStat(c *gin.Context) {
	data, err := stat.WeekStat(14)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, data)
}

// WeekRatio 占比
func (w ConsoleController) WeekRatio(c *gin.Context) {
	data, err := stat.WeekRatio()
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, data)
}

// TopProvider 提供商
func (w ConsoleController) TopProvider(c *gin.Context) {
	w.RenderSuccess(c, nil)
}

// TopNode 节点
func (w ConsoleController) TopNode(c *gin.Context) {
	w.RenderSuccess(c, nil)
}
