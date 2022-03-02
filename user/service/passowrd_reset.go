package service

import (
	"errors"
	baseModel "github.com/869413421/pg-service/common/pkg/model"
	string2 "github.com/869413421/pg-service/common/pkg/string"
	"github.com/869413421/pg-service/user/pkg/repo"
	"gorm.io/gorm"
	"time"
)

// PasswordResetServiceInterface 重置密码业务接口
type PasswordResetServiceInterface interface {
	Reset(token string) (string, error)
}

// PasswordResetService 重置密码业务
type PasswordResetService struct {
	UserRepo          repo.UserRepositoryInterface
	PasswordResetRepo repo.PasswordRestRepositoryInterface
}

// NewPasswordResetService 创建业务层
func NewPasswordResetService() PasswordResetServiceInterface {
	return &PasswordResetService{
		UserRepo:          repo.NewUserRepository(),
		PasswordResetRepo: repo.NewPasswordResetRepository(),
	}
}

// Reset 重置密码,返回新的密码
func (srv *PasswordResetService) Reset(token string) (string, error) {
	//1.根据token获取密码重置记录
	passwordReset, err := srv.PasswordResetRepo.GetByToken(token)
	if err != nil {
		return "", err
	}

	//2.比较时间，查看邮件是否已经超时或已验证
	if passwordReset.Verify == 1 {
		return "", errors.New("the record has been verified")
	}
	d, _ := time.ParseDuration("+5m")
	expTime := passwordReset.CreatedAt.Add(d)
	if time.Now().After(expTime) {
		return "", errors.New("verify that the message has expired")
	}

	//3.获取用户更新密码(使用事務)
	newPassword := string2.RandString(8)
	db := baseModel.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		user, err := srv.UserRepo.GetByEmail(passwordReset.Email)
		if err != nil {
			return err
		}
		user.Password = newPassword
		result := tx.Debug().Save(&user)
		err = result.Error
		if err != nil {
			return err
		}
		rowsAffected := result.RowsAffected
		if rowsAffected == 0 {
			return errors.New("update password fail")
		}

		//4.更新重置记录
		passwordReset.Verify = 1
		result = tx.Debug().Save(&passwordReset)
		err = result.Error
		if err != nil {
			return err
		}
		rowsAffected = result.RowsAffected
		if rowsAffected == 0 {
			return errors.New("update password fail")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return newPassword, nil
}
