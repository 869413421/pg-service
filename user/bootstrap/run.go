package bootstrap

import (
	baseModel "github.com/869413421/pg-service/common/pkg/model"
	"github.com/869413421/pg-service/common/pkg/rabbitmq"
	"github.com/869413421/pg-service/user/handler"
	"github.com/869413421/pg-service/user/pkg/model"
	pb "github.com/869413421/pg-service/user/proto/user"
	subscriber2 "github.com/869413421/pg-service/user/subscriber"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"
)

func Run() {
	//1.准备数据库连接，并且执行数据库迁移
	db := baseModel.GetDB()
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.PasswordReset{})

	//2.创建服务,初始化服务
	service := micro.NewService(
		micro.Name("pg.service.user"),
		micro.Version("v1"),
	)
	service.Init()

	//3.启动订阅
	brk, err := rabbitmq.GetBroker()
	if err != nil {
		log.Fatal("connection rabbitmq error", err)
	}
	eventSubscriber := subscriber2.NewEventSubscriber(brk)
	err = eventSubscriber.Subscriber()
	if err != nil {
		log.Fatal(err)
		return
	}

	//4.注册服务处理器
	err = pb.RegisterUserServiceHandler(service.Server(), handler.NewUserServiceHandler())
	if err != nil {
		return
	}

	//5.启动服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
