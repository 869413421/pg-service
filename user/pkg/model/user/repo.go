package user

import "github.com/869413421/pg-service/common/pkg/model"

type UserRepositoryInterface interface {
	GetByID(id uint64) (user User, err error)
}

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo UserRepository) GetByID(id uint64) (user User, err error) {
	err = model.DB.First(&user, id).Error
	return
}
