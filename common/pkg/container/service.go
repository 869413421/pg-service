package container

import (
	userPb "github.com/869413421/pg-service/user/proto/user"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
)

var service micro.Service
var userServiceClient userPb.UserService

// SetService 设置服务实例
func SetService(srv micro.Service) {
	service = srv
}

// GetService 返回服务实例
func GetService() micro.Service {
	return service
}

// GetServiceBroker 返回服务Broker实例
func GetServiceBroker() broker.Broker {
	return service.Options().Broker
}

// SetUserServiceClient 设置客户端实例
func SetUserServiceClient(userService userPb.UserService) {
	userServiceClient = userService
}

// GetUserServiceClient 获取客户端实例
func GetUserServiceClient() userPb.UserService {
	return userServiceClient
}
