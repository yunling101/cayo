package sender

import (
	"bytes"
	"fmt"
)

// FeiShu
type FeiShu struct{ *Channel }

// Sender
func (c *FeiShu) Delivery() (r Response) {
	b, _ := c.Marshal(map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title": "Cayo云监控告警触发通知",
					"content": [][]map[string]interface{}{
						{
							{"tag": "text", "text": fmt.Sprintf("告警ID: %v\n", c.Alert.ID)},
							{"tag": "text", "text": fmt.Sprintf("告警标题: %s\n", c.Alert.Title)},
							{"tag": "text", "text": fmt.Sprintf("任务名称: %s\n", c.Alert.TaskName)},
							{"tag": "text", "text": fmt.Sprintf("告警级别: %s\n", c.Alert.Level)},
							{"tag": "text", "text": fmt.Sprintf("告警阀值: %s\n", c.Alert.Threshold)},
							{"tag": "text", "text": fmt.Sprintf("当前值: %v\n", c.Alert.CurValue)},
							{"tag": "text", "text": fmt.Sprintf("指标名称: %s\n", c.Alert.Metric)},
							{
								"tag":  "text",
								"text": fmt.Sprintf("告警时间: %s", c.Alert.CreateTime.Format("2006-01-02 15:04:05")),
							},
						},
					},
				},
			},
		},
	})
	if s, err := c.Requests(c.Contact.WebHook, bytes.NewBuffer(b)); err != nil {
		r.Err = err
		return
	} else {
		var d map[string]interface{}
		if err := c.Unmarshal(s, &d); err == nil {
			if d["StatusCode"].(float64) != 0 {
				r.Err = fmt.Errorf("%s", d["StatusMessage"].(string))
				return
			}
		} else {
			r.Err = err
		}
	}
	return
}
