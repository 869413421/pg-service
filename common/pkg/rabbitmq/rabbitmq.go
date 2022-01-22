package rabbitmq

import (
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
		return nil, err
	}
	if err := b.Connect(); err != nil {
		return nil, err
	}
	return b, nil
}
