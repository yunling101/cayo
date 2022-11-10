package task

import (
	"time"

	"github.com/yunling101/cayo/pkg/global"
)

// Task 任务/站点
type Task struct {
	ID              int       `json:"id"`                                       // ID
	TaskId          string    `json:"task_id" bson:"task_id"`                   // 任务ID
	Name            string    `json:"name"`                                     // 任务名称
	Protocol        string    `json:"protocol"`                                 // 协议
	Status          bool      `json:"status"`                                   // 监控状态
	AlarmState      bool      `json:"alarm_state" bson:"alarm_state"`           // 报警状态
	Availability    int       `json:"availability"`                             // 可用率
	Address         string    `json:"address"`                                  // 地址，多个\n隔开
	PackNumber      int       `json:"pack_number" bson:"pack_number"`           // 包数，ping
	Frequency       int       `json:"frequency"`                                // 监控评率，分钟
	Port            int       `json:"port"`                                     // 端口
	Retry           int       `json:"retry"`                                    // 重试次数
	Method          string    `json:"method"`                                   // 请求方法
	RequestMode     string    `json:"request_mode" bson:"request_mode"`         // 请求方式
	RequestContent  string    `json:"request_content" bson:"request_content"`   // 请求内容
	ResponseMode    string    `json:"response_mode" bson:"response_mode"`       // 响应方式
	ResponseContent string    `json:"response_content" bson:"response_content"` // 响应内容
	ResponseTime    string    `json:"response_time" bson:"response_time"`       // 响应时间
	ResolutionType  string    `json:"resolution_type" bson:"resolution_type"`   // 解析类型
	Server          string    `json:"server"`                                   // 服务器
	Tls             bool      `json:"tls"`                                      // 安全连接
	Header          string    `json:"header"`                                   // HTTP请求头
	Cookie          string    `json:"cookie"`                                   // Cookie
	Username        string    `json:"username"`                                 // 用户
	Password        string    `json:"password"`                                 // 密码
	SSLVerify       bool      `json:"ssl_verify" bson:"ssl_verify"`             // 证书验证
	ProbeNode       []int     `json:"probe_node" bson:"probe_node"`             // 探测节点
	Alarm           []int     `json:"alarm"`                                    // 告警规则
	Notify          []int     `json:"notify"`                                   // 通知联系人
	CreateTime      time.Time `json:"create_time" bson:"create_time"`           // 添加时间
}

// Metric 指标数据
type Metric struct {
	ID           int       `json:"id"`                                 // ID
	TaskId       string    `json:"task_id" bson:"task_id"`             // 任务ID
	ProbePoint   string    `json:"probe_point" bson:"probe_point"`     // 探测节点
	ProbeSource  string    `json:"probe_source" bson:"probe_source"`   // 探测源
	Target       string    `json:"target"`                             // 探测目标
	ResponseTime float64   `json:"response_time" bson:"response_time"` // 响应时间
	ResponseCode int       `json:"response_code" bson:"response_code"` // 响应状态码
	ResponseBody string    `json:"response_body" bson:"response_body"` // 响应内容
	DnsServer    string    `json:"dns_server" bson:"dns_server"`       // DNS服务器
	DNSTime      float64   `json:"dns_time" bson:"dns_time"`           // DNS时间
	ConnectTime  float64   `json:"connect_time" bson:"connect_time"`   // 连接时间
	PackTime     float64   `json:"pack_time" bson:"pack_time"`         // 首包时间
	SSLTime      float64   `json:"ssl_time" bson:"ssl_time"`           // SSL证书解析时间
	SSLDays      int       `json:"ssl_days" bson:"ssl_days"`           // SSL证书到期天数
	PacketLoss   float64   `json:"packet_loss" bson:"packet_loss"`     // 包丢失的百分比
	Message      string    `json:"message"`                            // 错误信息
	StatusCode   int       `json:"status_code" bson:"status_code"`     // 探测状态 0: 异常 1: 正常
	CreateTime   time.Time `json:"create_time" bson:"create_time"`     // 添加时间
}

type MetricCache struct {
	ID         int       `json:"ID"`
	TaskId     string    `json:"task_id" bson:"task_id"` // 任务ID
	Value      int       `json:"value"`
	Status     bool      `json:"status"`
	Probe      []int     `json:"probe"`
	CreateTime time.Time `json:"create_time" bson:"create_time"` // 添加时间
}

// TableName 表名称
func (t Task) TableName() string {
	return global.Config().DataBase.Prefix + "task"
}

// TableName 表名称
func (t Metric) TableName() string {
	return global.Config().DataBase.Prefix + "metric_%s"
}

// TableName 表名称
func (t MetricCache) TableName() string {
	return global.Config().DataBase.Prefix + "metric_cache"
}

type TaskType int

const (
	Invalid TaskType = iota
	HTTP
	PING
	TCP
	UDP
	DNS
	SMTP
	// TRACEROUTE
)

var taskType = []string{
	Invalid: "Invalid",
	HTTP:    "HTTP",
	PING:    "PING",
	TCP:     "TCP",
	UDP:     "UDP",
	DNS:     "DNS",
	SMTP:    "SMTP",
}

func GetTaskType(name string) (p TaskType) {
	for st, n := range taskType {
		if n == name {
			return TaskType(st)
		}
	}
	return
}

func GetTaskTypeList() (names []*map[string]interface{}) {
	for st, n := range taskType {
		if n == "Invalid" {
			continue
		}
		names = append(names, &map[string]interface{}{
			"id":   st,
			"name": n,
		})
	}
	return
}

type DnsType int

const (
	TypeNone DnsType = iota
	TypeA
	TypeCNAME
	TypeTXT
	TypeMX
	TypeNS
)

var dnsType = []string{
	TypeNone:  "None",
	TypeA:     "A",
	TypeCNAME: "CNAME",
	TypeTXT:   "TXT",
	TypeMX:    "MX",
	TypeNS:    "NS",
}

func GetDnsType(name string) (p DnsType) {
	for st, n := range dnsType {
		if n == name {
			return DnsType(st)
		}
	}
	return
}

func GetDnsTypeList() (names []*map[string]interface{}) {
	for st, n := range dnsType {
		if n == "None" {
			continue
		}
		names = append(names, &map[string]interface{}{
			"id":   st,
			"name": n,
		})
	}
	return
}
