package main

import (
	handler "github.com/869413421/pg-service/user/handler"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("pg.service.user"),
		micro.Version("latest"),
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
