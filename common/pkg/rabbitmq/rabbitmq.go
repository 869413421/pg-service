package rabbitmq

import (
	"github.com/869413421/pg-service/common/pkg/logger"
	"github.com/go-acme/lego/v3/platform/config/env"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
)

// GetBroker 获取broker
func GetBroker() (broker.Broker, error) {
	b := rabbitmq.NewBroker(
		broker.Addrs(env.GetOrDefaultString("MICRO_BROKER_ADDRESS", "amqp://guest:guest@localhost:5672")),
	)
	if err := b.Init(); err != nil {
		logger.Danger("rabbitmq broker init error:",err)
		return nil, err
	}
	if err := b.Connect(); err != nil {
		logger.Danger("rabbitmq broker connect error:",err)
		return nil, err
	}
	return b, nil
}
