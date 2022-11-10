package user

import (
	"fmt"
	"time"

	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/toolkits/text"
)

func ValidUserPass(n alarm.AlarmContact) error {
	var user alarm.AlarmContact
	if err := q.Table(n.TableName()).QueryOne(
		q.M{"$or": []q.M{{"username": n.Username}}}, &user); err != nil {
		return fmt.Errorf(fmt.Sprintf("%s用户名不存在!", n.Username))
	}
	if !user.IsActive {
		return fmt.Errorf("%s", "此用户已被管理员停用!")
	}
	if user.Password != text.GenerateSha256ToString(n.Password) {
		return fmt.Errorf("%s", "密码错误!")
	}

	if err := q.Table(n.TableName()).UpdateOne(q.M{"username": user.Username}, q.M{"last_login": time.Now()}); err != nil {
		return fmt.Errorf("%s", "登录时间更新错误!")
	}
	return nil
}

func ModifyPassword(username, oldPassword, newPassword string) error {
	var user alarm.AlarmContact
	err := q.Table(user.TableName()).QueryOne(q.M{"username": username}, &user)
	if err != nil {
		return fmt.Errorf("%s用户不存在!", username)
	}
	if text.GenerateSha256ToString(oldPassword) != user.Password {
		return fmt.Errorf("%s!", "原密码错误")
	}
	err = q.Table(user.TableName()).UpdateOne(
		q.M{"username": username}, q.M{"password": text.GenerateSha256ToString(newPassword)})
	if err != nil {
		return fmt.Errorf("%s", "密码修改错误!")
	}
	return nil
}
