package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo/mongomgo"
	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/global"
)

func LoginSessions() gin.HandlerFunc {
	c := global.DB.C(global.Config().DataBase.Prefix + "session")
	store := mongomgo.NewStore(c, 3600, true, []byte("secret"))
	return sessions.Sessions("session_key", store)
}

func Login(c *gin.Context, username string) error {
	session := sessions.Default(c)

	session.Set("sessionid", username)
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("sessionid")

		if user == nil {
			data := make(map[string]interface{})
			data["code"] = 1001
			data["msg"] = "login"
			c.JSON(200, data)

			c.Abort()
			return
		}
		c.Next()
	}
}

func sessionDelete(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("sessionid")
	if user == nil {
		return
	}

	session.Delete("sessionid")
	if err := session.Save(); err != nil {
		return
	}
}

func LoginOut(c *gin.Context) {
	sessionDelete(c)
	c.Redirect(http.StatusFound, "/#/login")
	return
}
