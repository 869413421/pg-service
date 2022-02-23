package bootstrap

import (
	"github.com/869413421/pg-service/common/pkg/container"
	"github.com/869413421/pg-service/common/pkg/logger"
	baseModel "github.com/869413421/pg-service/common/pkg/model"
	pgPrometheus "github.com/869413421/pg-service/common/pkg/prometheus"
	"github.com/869413421/pg-service/common/pkg/trace"
	"github.com/869413421/pg-service/user/handler"
	"github.com/869413421/pg-service/user/pkg/model"
	pb "github.com/869413421/pg-service/user/proto/user"
	pgSubscriber "github.com/869413421/pg-service/user/subscriber"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimiter "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	tracePlugin "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"os"
	"time"
)

const QPS = 1000

func Run() {
	//1.准备数据库连接，并且执行数据库迁移
	db := baseModel.GetDB()
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.PasswordReset{})

	//2.初始化jaeger链路追踪
	t, io, err := trace.NewTracer("pg.user.service", os.Getenv("MICRO_TRACE_SERVER"))
	if err != nil {
		logger.Danger("init jaeger error:", err)
		return
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//3.初始化rpc服务
	service := micro.NewService(
		micro.Name("pg.service.user"),
		micro.Version("v1"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*20),
		micro.WrapHandler(tracePlugin.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		//令牌桶算法通过控制从外界发送到本节点的请求速率来限流
		micro.WrapHandler(ratelimiter.NewHandlerWrapper(QPS)),
	)
	service.Init()
	container.SetService(service)

	//3.启动订阅
	brk := service.Options().Broker
	err = brk.Init()
	if err != nil {
		logger.Danger("connection init error:", err)
		return
	}
	err = brk.Connect()
	if err != nil {
		logger.Danger("connection broker error:", err)
		return
	}
	defer brk.Disconnect()
	eventSubscriber := pgSubscriber.NewEventSubscriber(brk)
	err = eventSubscriber.Subscriber()
	if err != nil {
		logger.Danger("subscriber broker error:", err)
		return
	}

	//4.启动普罗米修斯
	pgPrometheus.PrometheusBoot()

	//4.注册服务处理器
	err = pb.RegisterUserServiceHandler(service.Server(), handler.NewUserServiceHandler())
	if err != nil {
		logger.Danger("register user service handler error:", err)
		return
	}

	//5.启动服务
	if err := service.Run(); err != nil {
		logger.Danger("service run error:", err)
	}
}
