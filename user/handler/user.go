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

// Create 创建用户
func (srv UserServiceHandler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	//1.验证提交信息
	user := model.User{}
	user.CreateFill(req)

	errs := requests.ValidateUserEdit(user)
	if len(errs) > 0 {
		errStr, _ := JsonHandler.Marshal(errs)
		return errors.BadRequest("User.Create.Validate.Error", string(errStr))
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

func (srv UserServiceHandler) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	//1.获取用户
	id := req.Id
	_user, err := srv.repo.GetByID(id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return errors.NotFound("User.Update.GetUserByID.Error", "user not found ,check you request id")
	}

	//2.验证提交信息
	_user.UpdateFill(req)
	errs := requests.ValidateUserEdit(*_user)
	if len(errs) > 0 {
		errStr, _ := JsonHandler.Marshal(errs)
		return errors.BadRequest("User.Update.Validate.Error", string(errStr))
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
