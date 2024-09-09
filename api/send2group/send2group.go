package send2group

import (
	"github.com/gin-gonic/gin"
	"github.com/morpingsss/ws-server/api"
	"github.com/morpingsss/ws-server/define/retcode"
	"github.com/morpingsss/ws-server/servers"
	"net/http"
)

type Controller struct{}

type inputData struct {
	SendUserId string `json:"sendUserId"`
	GroupName  string `json:"groupName" validate:"required"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Data       string `json:"data"`
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
	messageId := servers.SendMessage2Group(systemId, input.SendUserId, input.GroupName, input.Code, input.Msg, &input.Data)

	api.Render(ctx.Writer, retcode.SUCCESS, "success", map[string]string{
		"messageId": messageId,
	})
}
