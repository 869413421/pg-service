package main

import (
	"github.com/869413421/pg-service/common/pkg/logger"
	"github.com/869413421/pg-service/common/pkg/trace"
	"github.com/869413421/pg-service/common/pkg/wrapper/opentracing/gin2micro"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/web"
	"github.com/opentracing/opentracing-go"
	"os"
)

func main() {
	//1.创建web服务
	g := gin.Default()
	var serviceName = "pg.api.user"
	service := web.NewService(
		web.Name(serviceName),
		web.Address(":8888"),
		web.Handler(g),
	)

	// 2.初始化链路追踪
	gin2micro.SetSamplingFrequency(50)
	t, io, err := trace.NewTracer(serviceName, os.Getenv("MICRO_TRACE_SERVER"))
	if err != nil {
		logger.Danger("user api start tracer error:", err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 3.创建用户服务客户端
	cli := pb.NewUserService("pg.service.user", client.DefaultClient)

	v1 := g.Group("/user")
	v1.Use(gin2micro.TracerWrapper)
	{
		v1.GET("/get", func(context *gin.Context) {
			ctx, ok := gin2micro.ContextWithSpan(context)
			if ok == false {
				logger.Warning("user api user/get get context err")
			}
			req := &pb.GetRequest{}
			err := context.BindQuery(req)
			if err != nil {
				context.JSON(200, gin.H{
					"code": "500",
					"msg":  "bad request",
				})
			}
			if resp, err := cli.Get(ctx, req); err != nil {
				context.JSON(200, gin.H{
					"code": "500",
					"msg":  err.Error(),
				})
			} else {
				context.JSON(200, gin.H{
					"code": "200",
					"data": resp,
				})
			}
		})
	}

	err = service.Init()
	if err != nil {
		logger.Danger("Init api error:", err)
	}
	err = service.Run()
	if err != nil {
		logger.Danger("start api error:", err)
		return
	}
}
