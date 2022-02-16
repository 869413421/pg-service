package main

import (
	"github.com/869413421/pg-service/common/pkg/container"
	"github.com/869413421/pg-service/common/pkg/logger"
	"github.com/869413421/pg-service/common/pkg/trace"
	"github.com/869413421/pg-service/common/pkg/wrapper/opentracing/gin2micro"
	"github.com/869413421/pg-service/user-api/bootstarp"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/web"
	"github.com/opentracing/opentracing-go"
	"os"
	"time"
)

func main() {
	//1.创建gin并启动web服务
	g := bootstarp.SetupRoute()
	var serviceName = "pg.api.user"
	service := web.NewService(
		web.Name(serviceName),
		web.Address(":8888"),
		web.Handler(g),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*20),
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
	container.SetUserServiceClient(cli)
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
