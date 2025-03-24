package main

import (
	"cinexus/cmd"
	"cinexus/config"
	"cinexus/pkg/logger"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		panic("配置初始化失败: " + err.Error())
	}

	// 初始化日志
	if err := logger.Init(); err != nil {
		panic("日志初始化失败: " + err.Error())
	}
	defer logger.Sync()

	cmd.Execute()
}
