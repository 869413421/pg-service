package model

import (
	"github.com/869413421/pg-service/common/pkg/password"
	"github.com/jinzhu/gorm"
)

// BeforeSave 保存前模型事件
func (model *User) BeforeSave(tx *gorm.DB) (err error) {
	//1.如果密码没加密，进行一次加密
	if !password.IsHashed(model.Password) {
		model.Password, err = password.Hash(model.Password)
		if err!=nil{
			return err
		}
	}
	return nil
}