package register

import (
	"github.com/gin-gonic/gin"
	"github.com/morpingsss/ws-server/api"
	"github.com/morpingsss/ws-server/define/retcode"
	"github.com/morpingsss/ws-server/servers"
	"net/http"
)

type Controller struct{}

type inputData struct {
	SystemId string `json:"systemId" validate:"required"`
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

	err = servers.Register(input.SystemId)
	if err != nil {
		api.Render(ctx.Writer, retcode.FAIL, err.Error(), []string{})
		return
	}

	api.Render(ctx.Writer, retcode.SUCCESS, "success", []string{})
}
