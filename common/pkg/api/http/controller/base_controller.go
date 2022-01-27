package controller

import (
	"github.com/869413421/pg-service/common/pkg/message"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (*BaseController) ResponseJson(ctx *gin.Context, code int, errorMsg string, data interface{}) {
	responseData := message.ResponseData{
		Code:     code,
		ErrorMsg: errorMsg,
		Data:     data,
	}

	ctx.JSON(code, responseData)
	ctx.Abort()
}


