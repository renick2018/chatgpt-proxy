package main

import (
	"chatgpt-proxy/app"
	"chatgpt-proxy/component/chatgpt"
	"chatgpt-proxy/lib/logger"
)

func main() {
	chatgpt.LoadServers()
	// 加载多个APP的路由配置
	app.Include(app.Routers)
	// 初始化路由
	r := app.Init()
	if err := r.Run(":8088"); err != nil {
		logger.Error("start web listener err: ", err.Error())
	}
}
