package subscriber

import (
	"fmt"
	"github.com/micro/go-micro/v2/broker"
)

// 重置密码事件
const createPasswordResetTopic = "create.password.reset"

// EventSubscriberInterface 事件订阅者启动器接口
type EventSubscriberInterface interface {
	Subscriber() error
}

// EventSubscriber 事件订阅者启动器
type EventSubscriber struct {
	Broker broker.Broker
}

// NewEventSubscriber 创建事件订阅启动器
func NewEventSubscriber(brk broker.Broker) EventSubscriberInterface {
	return &EventSubscriber{Broker: brk}
}

// Subscriber 启动订阅
func (subscriber EventSubscriber) Subscriber() error {
	//1.重置密码事件订阅
	_, err := subscriber.Broker.Subscribe(createPasswordResetTopic, func(event broker.Event) error {
		// TODO 发送邮件
		fmt.Println("[sub] received message:", string(event.Message().Body), "header", event.Message().Header)
		return nil
	}, broker.Queue(createPasswordResetTopic), broker.DisableAutoAck())
	if err != nil {
		return err
	}

	return nil
}
