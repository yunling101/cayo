package sender

import (
	"bytes"
	"fmt"
)

type DingDing struct{ *Channel }

func (c *DingDing) Delivery() (r Response) {
	b, _ := c.Marshal(map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "Cayo云监控告警触发通知",
			"text": fmt.Sprintf(
				"**告警标题:** %s  \n **告警ID:** %v  \n **任务名称:** %s  \n **告警级别:** %s  \n **当前值:** %v  \n **告警时间:** %s",
				c.Alert.Title,
				c.Alert.ID,
				c.Alert.TaskName,
				c.Alert.Level,
				c.Alert.CurValue,
				c.Alert.CreateTime.Format("2006-01-02 15:04:05"),
			),
		},
	})
	if s, err := c.Requests(c.Contact.WebHook, bytes.NewBuffer(b)); err != nil {
		r.Err = err
		return
	} else {
		var d map[string]interface{}
		if err := c.Unmarshal(s, &d); err == nil {
			if d["errcode"].(float64) != 0 {
				r.Err = fmt.Errorf("%s", d["errmsg"].(string))
			}
		} else {
			r.Err = err
		}
	}
	return
}
