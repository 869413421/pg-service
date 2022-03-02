package main

import (
	"github.com/869413421/pg-service/common/pkg/container"
	"github.com/869413421/pg-service/common/pkg/logger"
	"github.com/869413421/pg-service/common/pkg/trace"
	"github.com/869413421/pg-service/common/pkg/wrapper/breaker/hystrix"
	"github.com/869413421/pg-service/common/pkg/wrapper/opentracing/gin2micro"
	"github.com/869413421/pg-service/user-api/bootstarp"
	_ "github.com/869413421/pg-service/user-api/docs"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/juju/ratelimit"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/web"
	ratelimiter "github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit/v2"
	"github.com/opentracing/opentracing-go"
	"os"
	"time"
)

const QPS = 1000

// @title 用户服务API
// @version 1.0
// @description 用户服务API

// @contact.name qingshui
// @contact.url https://qingshui.com
// @contact.email 13528685024@163.com

// @host localhost:8080
// @BasePath /user
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
	bucket := ratelimit.NewBucketWithRate(float64(QPS), int64(QPS))
	hystrix.Configure([]string{"pg.service.user.UserService.Auth"})
	clientService := micro.NewService(
		micro.Name("pg.api.user.cli"),
		micro.WrapClient(hystrix.NewClientWrapper()),
		//漏桶算法通过控制从本节点发出请求的速率来限流
		micro.WrapClient(ratelimiter.NewClientWrapper(bucket, false)),
	)

	// 4.创建用户服务客户端
	cli := pb.NewUserService("pg.service.user", clientService.Client())
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
