package routes

import (
	"github.com/869413421/pg-service/common/pkg/wrapper/opentracing/gin2micro"
	. "github.com/869413421/pg-service/user-api/app/http/controller"
	"github.com/gin-gonic/gin"
)

var userController = NewUserController()
var passwordController = NewPasswordResetController()

var middlewareHandlers []gin.HandlerFunc

//RegisterWebRoutes 注册路由
func RegisterWebRoutes(router *gin.Engine) {
	// 用户管理路由,所有路由必须包含user，因为micro网关只会映射路径中包含user的路由
	router.Use(gin2micro.TracerWrapper)
	api := router.Group("/")
	{
		api.POST("user/auth", userController.Auth)
	}
	{
		api.POST("user/password", passwordController.Store)
		api.PUT("user/password", passwordController.ResetPassword)
	}

	userApi := api.Group("user").Use(gin2micro.Jwt)
	{
		userApi.GET("", userController.Index)
		userApi.POST("", userController.Store)
		userApi.GET("/:id", userController.Show)
		userApi.PUT("/:id", userController.Update)
		userApi.DELETE("/:id", userController.Delete)
	}



}
