package types

import (
	"github.com/yunling101/cayo/pkg/propb"
	"google.golang.org/grpc"
)

// FlagConfig
type FlagConfig struct {
	Nid       int32  `json:"nid"`
	Server    string `json:"server"`
	Heartbeat int    `json:"heartbeat"`
	Interval  int    `json:"interval"`
	TimeOut   int    `json:"timeout"`
	Debug     bool   `json:"debug"`
}

// MemoryCache
type MemoryCache struct {
	Conn    *grpc.ClientConn    `json:"conn"`
	Rpc     propb.MonitorClient `json:"rpc"`
	Request *propb.Request      `json:"request"`
}

// Param
type Param struct {
	Data interface{} `json:"data"`
	Desc string      `json:"desc"`
}

// Assert
type Assert struct {
	Result bool `json:"result"`
	Index  int  `json:"index"`
}
