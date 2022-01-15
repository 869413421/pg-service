package handler

import (
	"context"
	"github.com/869413421/pg-service/user/pkg/model"
	"github.com/869413421/pg-service/user/pkg/repo"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/errors"
)

type UserServiceHandler struct {
	repo repo.UserRepositoryInterface
}

func NewUserServiceHandler() *UserServiceHandler {
	repo := repo.NewUserRepository()
	return &UserServiceHandler{repo: repo}
}

// GetByID 根据ID获取数据
func (srv *UserServiceHandler) GetByID(ctx context.Context, req *pb.GetByIDRequest, rsp *pb.GetByIDResponse) error {
	user, err := srv.repo.GetByID(req.GetId())
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return errors.BadRequest("User.GetByID", "user not found")
	}
	rsp.User = user.ToProtobuf()
	return nil
}

func (srv UserServiceHandler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	//1.检查邮箱是否重复
	user, err := srv.repo.GetByEmail(req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user.ID > 0 {
		return errors.BadRequest("User.Create", "email already exists ")
	}

	//2.检查电话是否重复
	user, err = srv.repo.GetByPhone(req.Phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user.ID > 0 {
		return errors.BadRequest("User.Create", "phone already exists ")
	}

	//3.创建用户
	user = &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RealName: req.RealName,
		Avatar:   req.Avatar,
		Phone:    req.Phone,
	}
	err = user.Store()
	if err != nil {
		return err
	}

	//4.返回用户信息
	rsp.User = user.ToProtobuf()
	return nil
}
