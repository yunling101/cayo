package sender

import (
	"crypto/tls"
	"fmt"

	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/toolkits/file"
	"gopkg.in/gomail.v2"
)

type MailObject struct{ *Channel }

// readFile
func (c *MailObject) readTemplate() (r []byte, err error) {
	template := fmt.Sprintf("%s/template/cayo.html", global.Config().Workspace)
	if !file.IsExist(template) {
		err = fmt.Errorf("%s", "模板不存在!")
		return
	}
	r, err = file.ReadAll(template)
	return
}

func (c *MailObject) Delivery() (r Response) {
	if !global.Config().Channel.Email.Enable {
		return
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", global.Config().Channel.Email.Username)
	msg.SetHeader("To", c.Contact.Email)
	msg.SetHeader("Subject", c.Alert.Title)

	if content, err := c.readTemplate(); err != nil {
		r.Err = err
		return
	} else {
		sendBody := fmt.Sprintf(string(content),
			c.Alert.Title,
			c.Alert.Title,
			fmt.Sprintf("事件: %s 名称: %s UUID: %s 指标: %s 阀值: %s 当前: %v",
				c.Alert.Title, c.Alert.TaskName, c.Alert.TaskID, c.Alert.Metric, c.Alert.Threshold, c.Alert.CurValue),
			c.Alert.Level,
			c.Alert.Metric,
			c.Alert.CreateTime.Format("2006-01-02 15:04:05"),
		)
		msg.SetBody("text/html", sendBody)
	}
	m := gomail.NewDialer(
		global.Config().Channel.Email.Host,
		global.Config().Channel.Email.Port,
		global.Config().Channel.Email.Username,
		global.Config().Channel.Email.Password,
	)
	m.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	r.Err = m.DialAndSend(msg)
	return
}
