package user

import (
	"github.com/869413421/pg-service/common/pkg/model"
)

type User struct {
	model.BaseModel
	Name     string `gorm:"column:name;type:varchar(255);not null;unique"`
	Email    string `gorm:"column:email;type:varchar(255) not null;unique"`
	Phone    string `gorm:"column:phone;type:varchar(255) not null;unique"`
	RealName string `gorm:"column:real_name;type:varchar(255);not null" `
	Avatar   string `gorm:"column:avatar;type:varchar(255);not null" `
	Status   string `gorm:"column:status;type:tinyint(1);default:0" `
}
