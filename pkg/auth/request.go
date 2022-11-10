package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/model/q"
)

type RequestParams struct {
	User alarm.AlarmContact
}

var Request RequestParams

func RequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("sessionid")

		if username != nil {
			err := q.Table(alarm.AlarmContact{}.TableName()).QueryOne(q.M{"username": username.(string)}, &Request.User)
			if err != nil {
				c.Abort()
				return
			}
		} else {
			// 用户注销后要清空
			Request = RequestParams{User: alarm.AlarmContact{}}
		}
		c.Next()
	}
}
