package handler

import (
	"context"
	. "github.com/869413421/pg-service/common/pkg/encoder"
	"github.com/869413421/pg-service/user/pkg/model"
	"github.com/869413421/pg-service/user/pkg/repo"
	"github.com/869413421/pg-service/user/pkg/requests"
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
	//1.验证提交信息
	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RealName: req.RealName,
		Avatar:   req.Avatar,
		Phone:    req.Phone,
	}

	errs := requests.ValidateUserEdit(user)
	if len(errs) > 0 {
		errStr, _ := JsonHandler.Marshal(errs)
		return errors.Unauthorized("User.Create", string(errStr))
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
