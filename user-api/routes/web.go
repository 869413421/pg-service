package routes

import (
	"github.com/869413421/pg-service/common/pkg/wrapper/opentracing/gin2micro"
	. "github.com/869413421/pg-service/user-api/app/http/controller"
	"github.com/gin-gonic/gin"
)

var userController = NewUserController()

var middlewareHandlers []gin.HandlerFunc

//RegisterWebRoutes 注册路由
func RegisterWebRoutes(router *gin.Engine) {
	// 用户管理路由
	router.Use(gin2micro.TracerWrapper)
	userApi := router.Group("/user")
	{
		userApi.GET("", userController.Index)
		userApi.POST("", userController.Store)
		userApi.GET("/:id", userController.Show)
		userApi.PUT("/:id", userController.Update)
		userApi.DELETE("/:id", userController.Delete)
	}
}
