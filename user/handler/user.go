package handler

import (
	"context"
	. "github.com/869413421/pg-service/common/pkg/encoder"
	"github.com/869413421/pg-service/common/pkg/password"
	"github.com/869413421/pg-service/common/pkg/types"
	"github.com/869413421/pg-service/user/pkg/model"
	"github.com/869413421/pg-service/user/pkg/repo"
	"github.com/869413421/pg-service/user/pkg/requests"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/869413421/pg-service/user/service"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserServiceHandler struct {
	UserRepo        repo.UserRepositoryInterface
	PasswordRepo    repo.PasswordRestRepositoryInterface
	TokenService    service.Authble
	PasswordService service.PasswordResetServiceInterface
}

func NewUserServiceHandler() *UserServiceHandler {
	userRepo := repo.NewUserRepository()
	passwordRepo := repo.NewPasswordResetRepository()
	tokenService := service.NewTokenService()
	passwordService := service.NewPasswordResetService()
	return &UserServiceHandler{UserRepo: userRepo, PasswordRepo: passwordRepo, TokenService: tokenService, PasswordService: passwordService}
}

// Get 根据ID获取数据
func (srv *UserServiceHandler) Get(ctx context.Context, req *pb.GetRequest, rsp *pb.UserResponse) error {
	user, err := srv.UserRepo.GetByID(req.GetId())
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return errors.BadRequest("User.GetByID", "user not found")
	}
	rsp.User = user.ToProtobuf()
	return nil
}

// Create 创建用户
func (srv *UserServiceHandler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.UserResponse) error {
	//1.验证提交信息
	user := &model.User{}
	types.Fill(user, req)
	errs := requests.ValidateUserEdit(*user)
	if len(errs) > 0 {
		errStr, _ := JsonHandler.Marshal(errs)
		return errors.Forbidden("User.Create.Validate.Error", string(errStr))
	}

	//2.创建用户
	err := user.Store()
	if err != nil {
		return err
	}

	//3.返回用户信息
	rsp.User = user.ToProtobuf()
	return nil
}

// Update 更新用户信息
func (srv *UserServiceHandler) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UserResponse) error {
	//1.获取用户
	id := req.Id
	_user, err := srv.UserRepo.GetByID(id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return errors.NotFound("User.Update.GetUserByID.Error", "user not found ,check you request id")
	}

	//2.验证提交信息
	types.Fill(_user, req)
	errs := requests.ValidateUserEdit(*_user)
	if len(errs) > 0 {
		errStr, _ := JsonHandler.Marshal(errs)
		return errors.Forbidden("User.Update.Validate.Error", string(errStr))
	}

	//3.更新用户
	rowsAffected, err := _user.Update()
	if rowsAffected == 0 || err != nil {
		return errors.InternalServerError("User.Update.Update.Error", err.Error())
	}

	//4.返回更新信息
	rsp.User = _user.ToProtobuf()
	return nil
}

// Delete 删除用户
func (srv *UserServiceHandler) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.UserResponse) error {
	//1.获取用户
	id := req.Id
	_user, err := srv.UserRepo.GetByID(id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return errors.NotFound("User.Delete.GetUserByID.Error", "user not found ,check you request id")
	}

	//2.删除用户
	rowsAffected, err := _user.Delete()
	if err != nil {
		return errors.InternalServerError("User.Delete.Delete.Error", err.Error())
	}
	if rowsAffected == 0 {
		return errors.BadRequest("User.Delete.Delete.Fail", "update fail")
	}

	//3.返回更新信息
	rsp.User = _user.ToProtobuf()
	return nil
}

// Auth 认证获取token
func (srv UserServiceHandler) Auth(ctx context.Context, req *pb.AuthRequest, rsp *pb.TokenResponse) error {
	//1.根据邮件获取用户
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := srv.UserRepo.GetByEmail(req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return errors.NotFound("User.Auth.GetByEmail.Error", "user not found ,check you request id")
	}

	//2.检验用户密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return errors.Unauthorized("User.Auth.CheckPassword.Error", err.Error())
	}

	//3.生成token
	token, err := srv.TokenService.Encode(user)
	if err != nil {
		return err
	}

	rsp.Token = token
	return nil
}

// ValidateToken 验证token
func (srv *UserServiceHandler) ValidateToken(ctx context.Context, req *pb.TokenRequest, rsp *pb.TokenResponse) error {
	claims, err := srv.TokenService.Decode(req.Token)
	if err != nil {
		return err
	}

	if claims.User.ID == 0 {
		return errors.BadRequest("User.ValidateToken.Error", "user invalid")
	}

	rsp.Valid = true

	return nil
}

//Pagination 分页
func (srv *UserServiceHandler) Pagination(ctx context.Context, req *pb.PaginationRequest, rsp *pb.PaginationResponse) error {
	users, pagerData, err := srv.UserRepo.Pagination(req.Page, req.PerPage)
	if err != nil {
		return errors.InternalServerError("user.Pagination.Pagination.Error", err.Error())
	}

	userItems := make([]*pb.User, len(users))
	for index, user := range users {
		userItem := user.ToProtobuf()
		userItems[index] = userItem
	}

	rsp.Users = userItems
	rsp.Total = pagerData.TotalCount
	return nil
}

//CreatePasswordReset 创建密码重置记录
func (srv *UserServiceHandler) CreatePasswordReset(ctx context.Context, req *pb.CreatePasswordResetRequest, rsp *pb.PasswordResetResponse) error {
	//1.获取提交邮箱,检查用户是否存在
	_, err := srv.UserRepo.GetByEmail(req.Email)
	if err != nil {
		return errors.NotFound("User.CreatePasswordReset.GetByEmail.Error", "user not found ,check you email")
	}

	passwordReset := &model.PasswordReset{}
	types.Fill(passwordReset, req)

	//2.生成md5保存到数据库
	passwordReset.Token = password.Md5Str(req.Email + time.Now().String())
	err = passwordReset.Store()
	if err != nil {
		return err
	}

	//err = srv.PasswordService.PushEvent(passwordReset)
	//if err != nil {
	//	return err
	//}

	//3.返回响应信息
	pbPasswordReset := &pb.PasswordReset{}
	types.Fill(pbPasswordReset, passwordReset)
	rsp.PasswordReset = pbPasswordReset

	return nil
}

// ResetPassword 重置密码
func (srv *UserServiceHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest, rsp *pb.ResetPasswordResponse) error {
	//1.执行重置逻辑
	newPassword, err := srv.PasswordService.Reset(req.Token)
	if err != nil {
		return err
	}

	//2.返回新密码
	rsp.ResetSuccess = true
	rsp.NewPassword = newPassword
	return nil
}
