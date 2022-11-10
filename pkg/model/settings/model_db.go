package settings

import (
	"time"

	"github.com/yunling101/cayo/pkg/global"
)

// Strategy
type Strategy struct {
	Name   string   `json:"name" bson:"name"`     // 名称
	Notify []Notify `json:"notify" bson:"notify"` // 通知
}

// Notify
type Notify struct {
	Name   string `json:"name" bson:"name"`     // 中文名称
	Alias  string `json:"alias" bson:"alias"`   // 英文名称
	Enable bool   `json:"enable" bson:"enable"` // 开启状态
}

// Settings
type Settings struct {
	Name           string     `json:"name" bson:"name"`                       // 固定名称
	UserLocal      bool       `json:"user_local" bson:"user_local"`           // 本地用户是否开启
	NotifyStrategy []Strategy `json:"notify_strategy" bson:"notify_strategy"` // 分派策略
	SendTime       bool       `json:"send_time" bson:"send_time"`             // 分派时间是否开启
	StartTime      string     `json:"start_time" bson:"start_time"`           // 开始时间
	EndTime        string     `json:"end_time" bson:"end_time"`               // 结束时间
	ModifyTime     time.Time  `json:"modify_time" bson:"modify_time"`         // 修改时间
}

// TableName
func (t Settings) TableName() string {
	return global.Config().DataBase.Prefix + "settings"
}
