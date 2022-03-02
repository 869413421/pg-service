package model

import (
	"fmt"
	"github.com/869413421/pg-service/common/pkg/container"
	"github.com/869413421/pg-service/common/pkg/encoder"
	"github.com/micro/go-micro/v2/broker"
	"gorm.io/gorm"
)

var createTopic = "create.password.reset"

func (model *PasswordReset) AfterCreate(tx *gorm.DB) (err error) {
	if model.Email != "" {
		err := pushCreateEvent(model)
		if err != nil {
			return err
		}
	}
	return nil
}

// pushCreateEvent 推送创建消息
func pushCreateEvent(model *PasswordReset) error {
	//1.获取发布连接
	publisher := container.GetServiceBroker()

	//2.构建broker消息
	body, err := encoder.JsonHandler.Marshal(model)
	if err != nil {
		return err
	}
	msg := &broker.Message{
		Header: map[string]string{
			"email": model.Email,
		},
		Body: body,
	}

	//3.发送消息到mq
	err = publisher.Publish(createTopic, msg)
	if err != nil {
		return err
	} else {
		fmt.Println("[pub] pubbed message:", string(msg.Body))
	}
	return nil
}
