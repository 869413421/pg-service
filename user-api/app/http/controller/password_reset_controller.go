package controller

import (
	"github.com/869413421/pg-service/common/pkg/api/http/controller"
	"github.com/869413421/pg-service/common/pkg/container"
	"github.com/869413421/pg-service/common/pkg/logger"
	"github.com/869413421/pg-service/common/pkg/wrapper/opentracing/gin2micro"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PasswordResetController struct {
	controller.BaseController
}

func NewPasswordResetController() *PasswordResetController {
	return &PasswordResetController{}
}

func (controller *PasswordResetController) Store(context *gin.Context) {
	//1.获取context中的信息
	ctx, ok := gin2micro.ContextWithSpan(context)
	if ok == false {
		logger.Warning("password_reset api user get context err")
	}

	//2.构建微服务请求体
	req := &pb.CreatePasswordResetRequest{}
	client := container.GetUserServiceClient()
	err := context.BindJSON(req)
	if err != nil {
		controller.ResponseJson(context, http.StatusForbidden, "json params error", []string{})
		return
	}

	//3.发起创建请求
	rsp, err := client.CreatePasswordReset(ctx, req)
	if err != nil {
		controller.ResponseJson(context, http.StatusInternalServerError, err.Error(), []string{})
		return
	}

	//4.响应信息
	controller.ResponseJson(context, http.StatusOK, "", rsp.PasswordReset)
}

func (controller *PasswordResetController) ResetPassword(context *gin.Context) {
	//1.获取context中的信息
	ctx, ok := gin2micro.ContextWithSpan(context)
	if ok == false {
		logger.Warning("password_reset api user get context err")
	}

	//2.构建微服务请求体
	req := &pb.ResetPasswordRequest{}
	client := container.GetUserServiceClient()
	err := context.BindJSON(req)
	if err != nil {
		controller.ResponseJson(context, http.StatusForbidden, "json params error", []string{})
		return
	}

	//3.发起创建请求
	rsp, err := client.ResetPassword(ctx, req)
	if err != nil {
		controller.ResponseJson(context, http.StatusInternalServerError, err.Error(), []string{})
		return
	}

	//4.响应信息
	controller.ResponseJson(context, http.StatusOK, "", rsp)
}


