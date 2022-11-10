package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/auth"
	"github.com/yunling101/cayo/pkg/cayoserver/controller"
	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/model/user"
	"github.com/yunling101/cayo/pkg/types"
)

// WebController
type WebController struct {
	controller.BaseController
	Bind struct {
		Web types.Web
	}
}

func (w WebController) Login(c *gin.Context) {
	if err := c.ShouldBindJSON(&w.Bind.Web.Login); err != nil {
		w.RenderFail(c, "用户或密码错误!")
		return
	}
	authView := alarm.AlarmContact{
		Username: w.Bind.Web.Login.Username,
		Password: w.Bind.Web.Login.Password,
	}
	if err := user.ValidUserPass(authView); err != nil {
		w.RenderFail(c, err.Error())
		return
	}
	err := auth.Login(c, w.Bind.Web.Login.Username)
	if err != nil {
		w.RenderFail(c, "登录失败, 未知的异常错误!")
		return
	}
	w.RenderSuccess(c, "/")
}

func (w WebController) LogOut(c *gin.Context) {
	auth.LoginOut(c)
}
