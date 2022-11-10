package settings

import (
	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/auth"
	"github.com/yunling101/cayo/pkg/cayoserver/controller"
	"github.com/yunling101/cayo/pkg/model/settings"
	"github.com/yunling101/cayo/pkg/model/user"
	"github.com/yunling101/cayo/pkg/types"
)

// SettingsController
type SettingsController struct {
	controller.BaseController
	Bind struct {
		Settings       settings.Settings
		ModifyPassword struct {
			OldPassword  string
			New1Password string
			New2Password string
		}
	}
}

// GetUserSettings
func (w SettingsController) GetUserSettings(c *gin.Context) {
	query := settings.Settings{}.GetUserSettings()
	w.RenderSuccess(c, query)
}

// ModifyUserSettings
func (w SettingsController) ModifyUserSettings(c *gin.Context) {
	if v := w.BindJSON(c, &w.Bind.Settings); !v {
		return
	}
	err := w.Bind.Settings.Save()
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

func (w SettingsController) ModifyUserPassword(c *gin.Context) {
	if v := w.BindJSON(c, &w.Bind.ModifyPassword); !v {
		return
	}
	if v := w.Params(c, []types.Param{
		{Data: w.Bind.ModifyPassword.OldPassword, Desc: "原密码"},
		{Data: w.Bind.ModifyPassword.New1Password, Desc: "新密码"},
		{Data: w.Bind.ModifyPassword.New2Password, Desc: "确认新密码"},
	}); !v {
		return
	}
	if w.Bind.ModifyPassword.OldPassword == w.Bind.ModifyPassword.New1Password {
		w.RenderFail(c, "原密码不能与新密码相同!")
		return
	}
	if w.Bind.ModifyPassword.New1Password != w.Bind.ModifyPassword.New2Password {
		w.RenderFail(c, "两次密码不一致!")
		return
	}
	if err := user.ModifyPassword(
		auth.Request.User.Username,
		w.Bind.ModifyPassword.OldPassword, w.Bind.ModifyPassword.New1Password); err != nil {
		w.RenderFail(c, err.Error())
		return
	}
	w.RenderSuccess(c, nil)
}
