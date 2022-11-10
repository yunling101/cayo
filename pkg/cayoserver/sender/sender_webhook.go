package sender

import (
	"bytes"
)

type WebHook struct{ *Channel }

func (c *WebHook) Delivery() (r Response) {
	s, _ := c.Marshal(map[string]interface{}{
		"id":            c.Alert.ID,
		"alarmInstance": "cayomonitor",
		"alertState":    c.Alert.Status,
		"curValue":      c.Alert.CurValue,
		"calcValue":     c.Alert.Threshold,
		"taskName":      c.Alert.TaskName,
		"taskId":        c.Alert.TaskID,
		"metricName":    c.Alert.Metric,
		"counter":       c.Alert.Counter,
		"triggerLevel":  c.Alert.Level,
		"createTime":    c.Alert.CreateTime,
		"timestamp":     c.Alert.CreateTime.Unix(),
	})
	if _, err := c.Requests(c.Contact.WebHook, bytes.NewBuffer(s)); err != nil {
		r.Err = err
	}
	return
}
