package controller

import (
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

// Index
// Anything godoc
// @Param token query string true "token"
// @Param page query int true "页码" default(1)
// @Param pageSize query int true "分页数量" default(10)
// @Summary 获取用户数据
// @Description 获取用户数据，支持分页
// @Accept  json
// @Produce  json
// @Success 200 {string} string "{"code":200,"errorMsg":"","data":{"users":[{"id":1,"password":"$2a$14$gJ6Iq2.cJ75v34OK.Mw/puJ9qZVcE79AESQa5AOBA6IzYbk/ukhxi","name":"huangyanming","email":"13528685024@163.com","phone":"13528685024","real_name":"黄彦铭11199","avatar":"121312312312","create_at":"2022-01-02 12:23:23","update_at":"2022-03-02 02:26:47"}],"total":1}}"
// @Router / [get]
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
	//1.获取context中的信息
	ctx, ok := gin2micro.ContextWithSpan(context)
	if ok == false {
		logger.Warning("user api user/store get context err")
	}

	//2.构建微服务请求体
	req := &pb.CreateRequest{}
	client := container.GetUserServiceClient()
	err := context.BindJSON(req)
	if err != nil {
		controller.ResponseJson(context, http.StatusForbidden, "json params error", []string{})
		return
	}

	//3.发起创建请求
	rsp, err := client.Create(ctx, req)
	if err != nil {
		controller.ResponseJson(context, http.StatusInternalServerError, err.Error(), []string{})
		return
	}

	//4.响应用户信息
	controller.ResponseJson(context, http.StatusOK, "", rsp.User)
}

func (controller *UserController) Update(context *gin.Context) {
	//1.获取context中的信息
	ctx, ok := gin2micro.ContextWithSpan(context)
	if ok == false {
		logger.Warning("user api user/update get context err")
	}

	//2.获取路由中的ID
	idStr := context.Param("id")
	if idStr == "" {
		controller.ResponseJson(context, http.StatusForbidden, "route id required", []string{})
		return
	}

	//3.构建微服务请求体
	req := &pb.UpdateRequest{}
	err := context.BindJSON(req)
	if err != nil {
		controller.ResponseJson(context, http.StatusForbidden, "json params error", []string{})
		return
	}
	id, _ := types.StringToInt(idStr)
	req.Id = uint64(id)

	//4.调用服务请求
	client := container.GetUserServiceClient()
	rsp, err := client.Update(ctx, req)
	if err != nil {
		controller.ResponseJson(context, http.StatusInternalServerError, err.Error(), []string{})
		return
	}

	//5.响应用户信息
	controller.ResponseJson(context, http.StatusOK, "", rsp.User)
}

func (controller *UserController) Delete(context *gin.Context) {
	//1.获取context中的信息
	ctx, ok := gin2micro.ContextWithSpan(context)
	if ok == false {
		logger.Warning("user api user/show get context err")
	}

	//2.获取路由中的ID
	idStr := context.Param("id")
	if idStr == "" {
		controller.ResponseJson(context, http.StatusForbidden, "route id required", []string{})
		return
	}

	//3.构建微服务请求体发起请求
	id, _ := types.StringToInt(idStr)
	req := &pb.DeleteRequest{Id: uint64(id)}
	client := container.GetUserServiceClient()
	rsp, err := client.Delete(ctx, req)
	if err != nil {
		controller.ResponseJson(context, http.StatusInternalServerError, err.Error(), []string{})
		return
	}

	//4.响应用户信息
	controller.ResponseJson(context, http.StatusOK, "", rsp.User)
}

func (controller *UserController) Show(context *gin.Context) {
	//1.获取context中的信息
	ctx, ok := gin2micro.ContextWithSpan(context)
	if ok == false {
		logger.Warning("user api user/show get context err")
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

func (controller *UserController) Auth(context *gin.Context) {
	//1.获取context中的信息
	ctx, ok := gin2micro.ContextWithSpan(context)
	if ok == false {
		logger.Warning("user api user/auth get context err")
	}

	//2.构建微服务请求体
	req := &pb.AuthRequest{}
	err := context.BindJSON(req)
	if err != nil {
		controller.ResponseJson(context, http.StatusForbidden, "json params error", []string{})
		return
	}

	//3.发起请求
	client := container.GetUserServiceClient()
	rsp, err := client.Auth(ctx, req)
	if err != nil {
		controller.ResponseJson(context, http.StatusInternalServerError, err.Error(), []string{})
		return
	}

	//3.响应用户信息
	controller.ResponseJson(context, http.StatusOK, "", rsp)
}
