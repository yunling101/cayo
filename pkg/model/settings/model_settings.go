package settings

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/cayo/pkg/model/alarm"
)

const (
	userSettings = "u_settings"
)

var defaultLevelNotify = map[string][]Notify{
	alarm.GetAlarmLevel(alarm.Critical): {
		{Name: "钉钉", Alias: "dingding", Enable: false},
		{Name: "飞书", Alias: "feishu", Enable: true},
		{Name: "邮件", Alias: "email", Enable: false},
		{Name: "WebHook", Alias: "webhook", Enable: true},
	},
	alarm.GetAlarmLevel(alarm.Warn): {
		{Name: "钉钉", Alias: "dingding", Enable: false},
		{Name: "飞书", Alias: "feishu", Enable: false},
		{Name: "邮件", Alias: "email", Enable: true},
		{Name: "WebHook", Alias: "webhook", Enable: false},
	},
	alarm.GetAlarmLevel(alarm.Info): {
		{Name: "钉钉", Alias: "dingding", Enable: false},
		{Name: "飞书", Alias: "feishu", Enable: false},
		{Name: "邮件", Alias: "email", Enable: false},
		{Name: "WebHook", Alias: "webhook", Enable: false},
	},
}

// GetUserSettings
func (s Settings) GetUserSettings() (r Settings) {
	err := global.DB.C(s.TableName()).Find(bson.M{"name": userSettings}).One(&r)
	if err != nil {
		for _, v := range alarm.GetAlarmLevelMap() {
			r.NotifyStrategy = append(r.NotifyStrategy, Strategy{
				Name:   v,
				Notify: defaultLevelNotify[v],
			})
		}
		r.UserLocal = true
		r.StartTime = "00:00"
		r.EndTime = "23:59"
		r.Name = userSettings
		return
	}
	return
}

// ModifySettings
func (s Settings) Save() (err error) {
	s.ModifyTime = time.Now()
	s.Name = userSettings
	err = global.DB.C(s.TableName()).Update(bson.M{"name": userSettings}, bson.M{"$set": s})
	if err == mgo.ErrNotFound {
		err = global.DB.C(s.TableName()).Insert(&s)
	}
	return
}
