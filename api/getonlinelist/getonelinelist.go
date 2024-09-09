package getonlinelist

import (
	"github.com/gin-gonic/gin"
	"github.com/morpingsss/ws-server/api"
	"github.com/morpingsss/ws-server/define/retcode"
	"github.com/morpingsss/ws-server/servers"
	"net/http"
)

type Controller struct{}

type inputData struct {
	GroupName string      `json:"groupName" validate:"required"`
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
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
	ret := servers.GetOnlineList(&systemId, &input.GroupName)

	// 确保返回了非空的响应
	if ret == nil {
		api.Render(ctx.Writer, retcode.FAIL, "No online users found", []string{})
		return
	}
	api.Render(ctx.Writer, retcode.SUCCESS, "success", ret)
}
