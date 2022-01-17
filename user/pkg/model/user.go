package model

import (
	baseModel "github.com/869413421/pg-service/common/pkg/model"
	pb "github.com/869413421/pg-service/user/proto/user"
)

// User 用户模型
type User struct {
	baseModel.BaseModel
	Name     string `gorm:"column:name;type:varchar(255);not null;unique" valid:"name"`
	Email    string `gorm:"column:email;type:varchar(255) not null;unique" valid:"email"`
	Phone    string `gorm:"column:phone;type:varchar(255) not null;unique" valid:"phone"`
	RealName string `gorm:"column:real_name;type:varchar(255);not null" valid:"real_name"`
	Avatar   string `gorm:"column:avatar;type:varchar(255);not null" valid:"avatar"`
	Status   string `gorm:"column:status;type:tinyint(1);default:0" `
	Password string `gorm:"column:password;type:varchar(255) not null;" valid:"password"`
}

// ToORM protobuf转换为orm
func ToORM(protoUser *pb.User) *User {
	user := &User{}
	user.ID = protoUser.Id
	user.Email = protoUser.Email
	user.Name = protoUser.Name
	user.Avatar = protoUser.Avatar
	user.Phone = protoUser.Phone
	user.RealName = protoUser.RealName
	user.Password = protoUser.Password
	return user
}

// ToProtobuf orm转换为protobuf
func (model *User) ToProtobuf() *pb.User {
	user := &pb.User{}
	user.Id = model.ID
	user.Email = model.Email
	user.Name = model.Name
	user.Avatar = model.Avatar
	user.CreateAt = model.CreatedAtDate()
	user.UpdateAt = model.UpdatedAtDate()
	user.Phone = model.Phone
	user.RealName = model.RealName
	user.Password = model.Password
	return user
}

// Store 创建用户
func (model *User) Store() (err error) {
	result := baseModel.GetDB().Create(&model)
	err = result.Error
	if err != nil {
		return err
	}
	return nil
}

// Update 更新用户
func (model *User) Update() (err error) {
	result := baseModel.GetDB().Save(&model)
	err = result.Error
	if err != nil {
		return err
	}
	return nil
}
