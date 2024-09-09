package bind2group

import (
	"github.com/gin-gonic/gin"
	"github.com/morpingsss/ws-server/api"
	"github.com/morpingsss/ws-server/define/retcode"
	"github.com/morpingsss/ws-server/servers"
	log "github.com/sirupsen/logrus"

	"net/http"
)

type Controller struct{}

type inputData struct {
	ClientId  string `json:"clientId" validate:"required"`
	GroupName string `json:"groupName" validate:"required"`
	UserId    string `json:"userId"`
	Extend    string `json:"extend"` // 拓展字段，方便业务存储数据
}

func (c *Controller) Run(ctx *gin.Context) {
	var input inputData
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := api.Validate(input)
	if err != nil {
		api.Render(ctx.Writer, retcode.FAIL, err.Error(), []string{})
		return
	}

	systemId := ctx.GetHeader("SystemId")
	log.Infoln("systemId: ", systemId)
	servers.AddClient2Group(systemId, input.GroupName, input.ClientId, input.UserId, input.Extend)
	log.Infoln("AddClient2Group: ", input.ClientId, input.GroupName, input.UserId, input.Extend)
	api.Render(ctx.Writer, retcode.SUCCESS, "success", []string{})
}
