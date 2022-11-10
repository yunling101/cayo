package sender

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/yunling101/cayo/pkg/model/alarm"
)

type Channel struct {
	Type    string             `json:"type"`
	Alert   alarm.AlarmTask    `json:"alert"`
	Contact alarm.AlarmContact `json:"contact"`
}

type Service interface {
	Delivery() Response
}

type Response struct {
	Data interface{}
	Err  error
}

func (c *Channel) NewChannel() Service {
	switch c.Type {
	case "opsalert":
		return &OpsAlert{c}
	case "dingding":
		return &DingDing{c}
	case "feishu":
		return &FeiShu{c}
	case "email":
		return &MailObject{c}
	case "webhook":
		return &WebHook{c}
	}
	return nil
}

func (c *Channel) Unmarshal(data []byte, v interface{}) error {
	jsonNew := jsoniter.ConfigCompatibleWithStandardLibrary
	return jsonNew.Unmarshal(data, v)
}

func (c *Channel) Marshal(val interface{}) (b []byte, err error) {
	jsonNew := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err = jsonNew.Marshal(val)
	return
}

// Requests
func (c *Channel) Requests(remoteURL string, params *bytes.Buffer) ([]byte, error) {
	var result []byte

	client := &http.Client{}
	request, err := http.NewRequest("POST", remoteURL, params)
	if err != nil {
		return result, err
	}
	request.Header.Set("Content-Type", "application/json")
	if response, err := client.Do(request); err != nil {
		return result, err
	} else {
		if response.StatusCode != 200 {
			return result, errors.New("status code is not 200")
		}
		return ioutil.ReadAll(response.Body)
	}
}
