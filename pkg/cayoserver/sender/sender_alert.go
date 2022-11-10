package sender

import (
	"bytes"
	"fmt"
)

type OpsAlert struct{ *Channel }

func (c *OpsAlert) Delivery() (r Response) {
	s, _ := c.Marshal(c.Alert)
	if b, err := c.Requests(c.Contact.WebHook, bytes.NewBuffer(s)); err != nil {
		r.Err = err
		return
	} else {
		var d map[string]interface{}
		if err := c.Unmarshal(b, &d); err != nil {
			r.Err = err
			return
		}
		if d["code"].(float64) == 0 {
			r.Err = fmt.Errorf("%s", d["msg"].(string))
		}
	}
	return
}
