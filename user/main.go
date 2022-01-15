package main

import (
	baseModel "github.com/869413421/pg-service/common/pkg/model"
	"github.com/869413421/pg-service/user/handler"
	"github.com/869413421/pg-service/user/pkg/model"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

func main() {
	//1.准备数据库连接，并且执行数据库迁移
	db := baseModel.GetDB()
	db.AutoMigrate(&model.User{})

	// New Service
	service := micro.NewService(
		micro.Name("pg.service.user"),
		micro.Version("v1"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	pb.RegisterUserServiceHandler(service.Server(), handler.NewUserServiceHandler())

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
