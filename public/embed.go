package public

import (
	"embed"
	"io/fs"
	"net/http"
)

// 双斜线和go:embed 之前不能有空格
// 只能用在包一级的变量中，不能用在函数或方法中

//go:embed build/*
var FS embed.FS

func StaticFS(relativePath string) (string, http.FileSystem) {
	sub, _ := fs.Sub(FS, "build"+relativePath)
	return relativePath, http.FS(sub)
}
