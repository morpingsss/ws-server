package closeclient

import (
	"github.com/gin-gonic/gin"
	"github.com/morpingsss/ws-server/api"
	"github.com/morpingsss/ws-server/define/retcode"
	"github.com/morpingsss/ws-server/servers"
	"net/http"
)

type Controller struct{}

type inputData struct {
	ClientId string `json:"clientId" validate:"required"`
}

func (c *Controller) Run(ctx *gin.Context) {
	var input inputData
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := api.Validate(input); err != nil {
		api.Render(ctx.Writer, retcode.FAIL, err.Error(), []string{})
		return
	}

	systemId := ctx.GetHeader("SystemId")

	// 发送信息
	servers.CloseClient(input.ClientId, systemId)

	api.Render(ctx.Writer, retcode.SUCCESS, "success", map[string]string{})
}
