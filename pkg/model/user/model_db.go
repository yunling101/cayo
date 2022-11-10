package user

import "github.com/yunling101/cayo/pkg/model/alarm"

func InitUserAdmin() {
	u := alarm.AlarmContact{
		Channel:     "local",
		Username:    "admin",
		Nickname:    "超级管理员",
		Password:    "123123",
		Email:       "yunling101@gmail.com",
		Phone:       "18600000000",
		Role:        "admin",
		Group:       "default",
		IsActive:    true,
		IsSuperuser: true,
	}
	u.Add(u)
}
