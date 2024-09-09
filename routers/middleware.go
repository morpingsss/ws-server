package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/morpingsss/ws-server/api"
	"github.com/morpingsss/ws-server/define"
	"github.com/morpingsss/ws-server/define/retcode"
	"github.com/morpingsss/ws-server/pkg/etcd"
	"github.com/morpingsss/ws-server/servers"
	"github.com/morpingsss/ws-server/tools/util"
	"net/http"
)

// AccessTokenMiddleware 检查和验证 SystemId
func AccessTokenMiddleware(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		c.Abort()
		return
	}

	systemId := c.GetHeader("SystemId")

	if systemId == "" {
		api.Render(c.Writer, retcode.FAIL, "系统ID不能为空", []string{})
		c.Abort()
		return
	}

	if util.IsCluster() {
		resp, err := etcd.Get(define.ETCD_PREFIX_ACCOUNT_INFO + systemId)
		if err != nil || resp.Count == 0 {
			api.Render(c.Writer, retcode.FAIL, "系统ID无效", []string{})
			c.Abort()
			return
		}
	} else {
		if _, ok := servers.SystemMap.Load(systemId); !ok {
			api.Render(c.Writer, retcode.FAIL, "系统ID无效", []string{})
			c.Abort()
			return
		}
	}

	c.Next() // 使用 c.Next() 继续执行处理链上的下一个中间件或路由处理函数
}
