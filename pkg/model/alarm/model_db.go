package alarm

import (
	"time"

	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/cayo/pkg/types"
)

// AlarmRule 报警规则
type AlarmRule struct {
	types.Currency `bson:"-"`
	ID             int           `json:"id"`                             // ID
	Team           int           `json:"team"`                           // 团队ID
	Name           string        `json:"name"`                           // 规则名称
	Task           types.Model   `json:"task"`                           // 关联任务
	State          bool          `json:"state"`                          // 是否开启
	Status         string        `json:"status"`                         // 运行状态
	Notify         []int         `json:"notify"`                         // 通知对象
	Condition      string        `json:"condition"`                      // 规则条件
	Metric         string        `json:"metric"`                         // 指标名称
	Alarm          []types.Alarm `json:"alarm"`                          // 报警详情
	CreateTime     time.Time     `json:"create_time" bson:"create_time"` // 创建时间
}

// AlarmContact 联系人
type AlarmContact struct {
	ID          int       `json:"id"`                               // ID
	Username    string    `json:"username"`                         // 用户名
	Nickname    string    `json:"nickname"`                         // 姓名
	Password    string    `json:"password"`                         // 密码
	Channel     string    `json:"channel" binding:"required"`       // 通道
	Email       string    `json:"email"`                            // 邮箱
	Phone       string    `json:"phone"`                            // 手机号
	IsActive    bool      `json:"is_active" bson:"is_active"`       // 是否活动
	IsSuperuser bool      `json:"is_superuser" bson:"is_superuser"` // 是否超级管理员
	Role        string    `json:"role"`                             // 用户角色; admin OR user
	Group       string    `json:"group"`                            // 用户组
	DingDing    string    `json:"dingding"`                         // dingding
	Feishu      string    `json:"feishu"`                           // feishu
	WebHook     string    `json:"webhook"`                          // WebHook
	LastLogin   time.Time `json:"last_login" bson:"last_login"`     // 登录时间
	CreateTime  time.Time `json:"create_time" bson:"create_time"`   // 创建时间
}

// AlarmTask 报警表
type AlarmTask struct {
	ID           int           `json:"id"`                                 // ID
	TaskName     string        `json:"task_name" bson:"task_name"`         // 任务名称
	TaskID       string        `json:"task_id" bson:"task_id"`             // 任务ID
	RuleID       int           `json:"rule_id" bson:"rule_id"`             // 规则ID
	Users        []types.Users `json:"users" bson:"users"`                 // 通知用户
	NodeID       int           `json:"node_id" bson:"node_id"`             // 节点ID
	Team         int           `json:"team"`                               // 团队ID
	Level        string        `json:"level"`                              // 级别
	Title        string        `json:"title"`                              // 标题
	CurValue     float64       `json:"cur_value" bson:"cur_value"`         // 当前值
	Threshold    string        `json:"threshold"`                          // 报警阀值
	Unit         string        `json:"unit"`                               // 阀值单位
	Metric       string        `json:"metric"`                             // 指标名称
	Status       string        `json:"status"`                             // 报警状态；trigger OR recovery
	Counter      int           `json:"counter"`                            // 计数器
	Archive      bool          `json:"archive"`                            // 是否归档
	ArchiveTime  time.Time     `json:"archive_time" bson:"archive_time"`   // 归档时间
	RecoveryTime time.Time     `json:"recovery_time" bson:"recovery_time"` // 恢复时间
	CreateTime   time.Time     `json:"create_time" bson:"create_time"`     // 报警时间
}

// TableName 表名称
func (r AlarmRule) TableName() string {
	return global.Config().DataBase.Prefix + "alarm_rule"
}

// TableName 表名称
func (t AlarmContact) TableName() string {
	return global.Config().DataBase.Prefix + "alarm_contact"
}

// TableName 表名称
func (t AlarmTask) TableName() string {
	return global.Config().DataBase.Prefix + "alarm_task"
}

type (
	Metric     string
	AlarmLevel int
)

const (
	Unknown AlarmLevel = iota
	Critical
	Warn
	Info
)

var (
	AlarmMetricList    []map[string]string
	notifyChannel      []map[string]string
	localNotifyChannel []map[string]string
	alarmLevel         map[AlarmLevel]string
)

func init() {
	alarmLevel = map[AlarmLevel]string{
		Critical: "Critical",
		Warn:     "Warn",
		Info:     "Info",
	}
	AlarmMetricList = []map[string]string{
		{"name": "响应时间", "id": "ResponseTime", "unit": "ms"},
		{"name": "响应状态码", "id": "ResponseCode", "unit": ""},
		{"name": "可用探测点百分比", "id": "AvailablePercent", "unit": "%"},
		{"name": "可用探测点数量", "id": "AvailablePoint", "unit": ""},
	}
	notifyChannel = []map[string]string{
		{"name": "告警平台(OpsAlert)", "id": "opsalert"},
		{"name": "本地联系人", "id": "local"},
	}
	localNotifyChannel = []map[string]string{
		{"name": "钉钉", "id": "dingding"},
		{"name": "飞书", "id": "feishu"},
		{"name": "邮件", "id": "email"},
		{"name": "WebHook", "id": "webhook"},
	}
}

func GetMetricName(id string) string {
	for i := 0; i < len(AlarmMetricList); i++ {
		v := AlarmMetricList[i]
		if v["id"] == id {
			return v["name"]
		}
	}
	return "Unknown"
}

func GetMetricUnit(id string) string {
	for i := 0; i < len(AlarmMetricList); i++ {
		v := AlarmMetricList[i]
		if v["id"] == id {
			return v["unit"]
		}
	}
	return ""
}

func GetNotifyName(id string) string {
	for i := 0; i < len(notifyChannel); i++ {
		v := notifyChannel[i]
		if v["id"] == id {
			return v["name"]
		}
	}
	return "Unknown"
}

func GetLocalNotifyName(id string) string {
	for i := 0; i < len(localNotifyChannel); i++ {
		v := localNotifyChannel[i]
		if v["id"] == id {
			return v["name"]
		}
	}
	return "Unknown"
}

func GetNotifyChannel() []map[string]string {
	return notifyChannel
}

func GetLocalNotifyChannel() []map[string]string {
	return localNotifyChannel
}

func GetAlarmLevelMap() map[AlarmLevel]string {
	return alarmLevel
}

func GetAlarmLevel(k AlarmLevel) string {
	if v, ok := alarmLevel[k]; ok {
		return v
	}
	return "Unknown"
}
