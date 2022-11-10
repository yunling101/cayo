package g

import (
	"github.com/yunling101/toolkits/logger"
)

type logging struct {
	App   logger.Logger
	Error logger.Logger
}

// Logger 实例化全局
var Logger logging

// InitLogger 初始化配置
func InitLogger() {
	// Logger.App = logger.Logger{
	// 	AppPath: fmt.Sprintf("%s/logs/cayo.log", Config().Dir),
	// }
	// Logger.Error = logger.Logger{
	// 	AppPath: fmt.Sprintf("%s/logs/error.log", Config().Dir),
	// }
	Logger.App.Init()
	Logger.Error.Init()
}
