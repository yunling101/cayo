package types

// Contact 联系人
type Contact struct {
	Add addRequestBody
}

type addRequestBody struct {
	Version int `json:"version" binding:"required"`
}

type Alarm struct {
	Level      string `json:"level"`         // 报警级别
	Threshold  string `json:"threshold"`     // 报警阀值
	Continuity int    `json:"continuity"`    // 连续周期
	Unit       string `json:"unit" bson:"-"` // 指标单位
}

type Users struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
}

type NotifyRatio struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type AlarmIndex struct {
	ResponseTime float64
	StatusCode   int
	ResponseCode int
}
