package main

import (
	"github.com/869413421/pg-service/user/handler"
	pb "github.com/869413421/pg-service/user/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("pg.service.user"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterUserServiceHandler(srv.Server(), handler.NewUserServiceHandler())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
