package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/morpingsss/ws-server/api"
	"github.com/morpingsss/ws-server/define/retcode"
	"github.com/morpingsss/ws-server/tools/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

type Controller struct {
}

type renderData struct {
	ClientId string `json:"clientId"`
}

func (controller *Controller) Run(ctx *gin.Context) {
	log.Infoln("Run start :")
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Errorf("upgrade error: %v", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer conn.Close()

	// 设置读取消息大小上限
	conn.SetReadLimit(maxMessageSize)

	// 解析参数
	systemId := ctx.Query("systemId")
	if len(systemId) == 0 {
		_ = Render(conn, "", "", retcode.SYSTEM_ID_ERROR, "系统ID不能为空", []string{})
		_ = conn.Close()
		return
	}

	clientId := util.GenClientId()
	clientSocket := NewClient(clientId, systemId, conn)

	// 将客户端添加到管理器
	Manager.AddClient2SystemClient(systemId, clientSocket)

	// 读取客户端消息
	clientSocket.Read()

	// 返回连接成功信息
	if err = api.ConnRender(conn, renderData{ClientId: clientId}); err != nil {
		_ = conn.Close()
		return
	}

	// 用户连接事件
	Manager.Connect <- clientSocket
}
