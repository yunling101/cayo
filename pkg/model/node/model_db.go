package node

import (
	"time"

	"github.com/yunling101/cayo/pkg/global"
)

// EndPoint 节点
type EndPoint struct {
	ID         int       `json:"id"`                          // ID
	Attribute  string    `json:"attribute"`                   // 归属
	Operator   string    `json:"operator" binding:"required"` // 运营商
	OperatorZh string    `json:"operator_zh" bson:"-"`
	TasksCount int       `json:"tasks_count" bson:"-"`
	Province   string    `json:"province"`                       // 省
	City       string    `json:"city"`                           // 市
	IP         string    `json:"ip"`                             // IP
	Status     string    `json:"status"`                         // 状态
	Selected   bool      `json:"selected"`                       // 默认选中状态
	Heartbeat  time.Time `json:"heartbeat"`                      // 心跳时间
	HostName   string    `json:"hostname"`                       // 主机名称
	Version    string    `json:"version"`                        // 版本
	Comment    string    `json:"comment"`                        // 备注
	CreateTime time.Time `json:"create_time" bson:"create_time"` // 添加时间
}

// TableName
func (t EndPoint) TableName() string {
	return global.Config().DataBase.Prefix + "end_point"
}

var operatorChannel = []map[string]string{
	{"id": "cmcc", "name": "中国移动"},
	{"id": "cdma", "name": "中国电信"},
	{"id": "unicom", "name": "中国联通"},
	{"id": "aliyun", "name": "阿里云"},
	{"id": "tencent", "name": "腾讯云"},
	{"id": "baidu", "name": "百度云"},
	{"id": "huawei", "name": "华为云"},
	{"id": "aws", "name": "亚马逊"},
	{"id": "google", "name": "谷歌云"},
	{"id": "other", "name": "其他"},
}

func GetOperatorMap() []map[string]string {
	return operatorChannel
}

func GetOperatorName(id string) string {
	for i := 0; i < len(operatorChannel); i++ {
		v := operatorChannel[i]
		if v["id"] == id {
			return v["name"]
		}
	}
	return "Unknown"
}
