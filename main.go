package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/morpingsss/ws-server/pkg/setting"
	"github.com/morpingsss/ws-server/routers"
	"github.com/morpingsss/ws-server/servers"
	"github.com/morpingsss/ws-server/tools/log"
	"github.com/morpingsss/ws-server/tools/util"
)

func main() {
	setting.Setup()
	log.Setup()

	router := gin.Default()

	// 初始化 API 和 WebSocket 路由
	initRoutes(router)

	// 如果是集群模式，初始化 RPC 服务
	if util.IsCluster() {
		servers.InitGRpcServer()
		fmt.Printf("启动RPC，端口号：%s\n", setting.CommonSetting.RPCPort)
	}

	// 启动 Manager 以处理客户端连接
	servers.StartManager()

	// 启动定时任务和消息发送
	go servers.WriteMessage()

	// 启动 HTTP 服务
	fmt.Printf("服务器启动成功，端口号：%s\n", setting.CommonSetting.HttpPort)
	if err := router.Run(":" + setting.CommonSetting.HttpPort); err != nil {
		panic(err)
	}
}

func initRoutes(router *gin.Engine) {
	// 初始化 API 路由
	routers.InitRoutes(router)
}
