package controller

import (
	"fmt"
	"github.com/869413421/pg-service/common/pkg/api/http/controller"
	"github.com/869413421/pg-service/common/pkg/container"
	"github.com/869413421/pg-service/common/pkg/logger"
	"github.com/869413421/pg-service/common/pkg/message"
	"github.com/869413421/pg-service/common/pkg/types"
	"github.com/869413421/pg-service/common/pkg/wrapper/opentracing/gin2micro"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	controller.BaseController
}

func NewUserController() *UserController {
	return &UserController{}
}

func (controller *UserController) Index(context *gin.Context) {
	//1.获取context中的信息
	ctx, ok := gin2micro.ContextWithSpan(context)
	if ok == false {
		logger.Warning("user api user get context err")
	}

	//2.构建参数发起请求
	pagination := &message.Pagination{}
	err := context.BindQuery(pagination)
	if err != nil {
		controller.ResponseJson(context, http.StatusForbidden, "pagination params required", []string{})
	}
	fmt.Println(1)
	fmt.Println(pagination)

	req := &pb.PaginationRequest{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	}
	client := container.GetUserServiceClient()
	rsp, err := client.Pagination(ctx, req)
	if err != nil {
		controller.ResponseJson(context, http.StatusInternalServerError, err.Error(), []string{})
		return
	}

	//4.响应用户信息
	controller.ResponseJson(context, http.StatusOK, "", rsp)
}

func (controller *UserController) Store(context *gin.Context) {

}

func (controller *UserController) Update(context *gin.Context) {

}

func (controller *UserController) Delete(context *gin.Context) {

}

func (controller *UserController) Show(context *gin.Context) {
	//1.获取context中的信息
	ctx, ok := gin2micro.ContextWithSpan(context)
	if ok == false {
		logger.Warning("user api user/get get context err")
	}

	//2.获取路由中的ID
	idStr := context.Param("id")
	if idStr == "" {
		controller.ResponseJson(context, http.StatusForbidden, "route id required", []string{})
		return
	}

	//3.构建微服务请求体发起请求
	id, _ := types.StringToInt(idStr)
	req := &pb.GetRequest{Id: uint64(id)}
	client := container.GetUserServiceClient()
	rsp, err := client.Get(ctx, req)
	if err != nil {
		controller.ResponseJson(context, http.StatusInternalServerError, err.Error(), []string{})
		return
	}

	//4.响应用户信息
	controller.ResponseJson(context, http.StatusOK, "", rsp.User)
}
