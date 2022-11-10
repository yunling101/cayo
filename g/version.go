package g

import (
	"fmt"
	"runtime"
)

var (
	// Version 版本号
	Version = "v1.0"

	// Branch 分支
	Branch string

	// BuildDate 构建日期
	BuildDate string

	// 构建hash
	CCGitHash string

	// GoVersion 构建版本
	GoVersion = runtime.Version()
)

// PrintVersion
func PrintVersion() string {
	return fmt.Sprintf("%s (goVersion: %s)", Version, GoVersion)
}

// Info
func Info() string {
	return fmt.Sprintf("(version=%s, branch=%s)", Version, Branch)
}

// PrintInfo
func PrintInfo() string {
	return fmt.Sprintf("(version=%s, branch=%s, buildDate=%s, gitHash=%s)",
		Version, Branch, BuildDate, CCGitHash)
}
