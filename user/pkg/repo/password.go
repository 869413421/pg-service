package repo

import (
	baseModel "github.com/869413421/pg-service/common/pkg/model"
	"github.com/869413421/pg-service/user/pkg/model"
	"gorm.io/gorm"
)

//PasswordRestRepositoryInterface 重置记录操作接口
type PasswordRestRepositoryInterface interface {
	GetByEmail(email string) (*model.PasswordReset, error)
	GetByToken(token string) (*model.PasswordReset, error)
}

//PasswordResetRepository 重置记录操作仓库
type PasswordResetRepository struct {
	DB *gorm.DB
}

// NewPasswordResetRepository 新建仓库
func NewPasswordResetRepository() PasswordRestRepositoryInterface {
	return &PasswordResetRepository{DB: baseModel.GetDB()}
}

// GetByEmail 根据邮件获取
func (repo *PasswordResetRepository) GetByEmail(email string) (*model.PasswordReset, error) {
	passwordReset := &model.PasswordReset{}
	err := repo.DB.Where("email = ?", email).First(&passwordReset).Error
	if err != nil {
		return nil, err
	}
	return passwordReset, nil
}

// GetByToken 根据token获取
func (repo *PasswordResetRepository) GetByToken(token string) (*model.PasswordReset, error) {
	passwordReset := &model.PasswordReset{}
	err := repo.DB.Where("token = ?", token).First(&passwordReset).Error
	if err != nil {
		return nil, err
	}
	return passwordReset, nil
}
