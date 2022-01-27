package container

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
)

var service micro.Service

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
