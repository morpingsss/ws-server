package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/morpingsss/ws-server/api/bind2group"
	"github.com/morpingsss/ws-server/api/closeclient"
	"github.com/morpingsss/ws-server/api/getonlinelist"
	"github.com/morpingsss/ws-server/api/register"
	"github.com/morpingsss/ws-server/api/send2client"
	"github.com/morpingsss/ws-server/api/send2clients"
	"github.com/morpingsss/ws-server/api/send2group"
	"github.com/morpingsss/ws-server/servers"
)

func InitRoutes(router *gin.Engine) {
	// 不需要中间件的路由组
	noAuthGroup := router.Group("/api")
	{
		noAuthGroup.POST("/register", (&register.Controller{}).Run)
	}

	// 需要中间件的路由组
	authGroup := router.Group("/api")
	{
		authGroup.Use(AccessTokenMiddleware) // 为该组应用中间件
		{
			authGroup.POST("/send_to_client", (&send2client.Controller{}).Run)
			authGroup.POST("/send_to_clients", (&send2clients.Controller{}).Run)
			authGroup.POST("/send_to_group", (&send2group.Controller{}).Run)
			authGroup.POST("/bind_to_group", (&bind2group.Controller{}).Run)
			authGroup.POST("/get_online_list", (&getonlinelist.Controller{}).Run)
			authGroup.POST("/close_client", (&closeclient.Controller{}).Run)
		}
	}

	// WebSocket 路由
	router.GET("/ws", (&servers.Controller{}).Run)
}
