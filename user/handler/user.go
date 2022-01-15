package handler

import (
	"context"
	"github.com/869413421/pg-service/user/pkg/model/user"
	pb "github.com/869413421/pg-service/user/proto/user"
)

type UserServiceHandler struct {
	repo user.UserRepositoryInterface
}

func NewUserServiceHandler() *UserServiceHandler {
	repo := user.NewUserRepository()
	return &UserServiceHandler{repo: repo}
}

func (srv *UserServiceHandler) GetByID(ctx context.Context, req *pb.GetByIDRequest, rsp *pb.GetByIDResponse) error {
	user, _ := srv.repo.GetByID(req.GetId())
	pbUser := &pb.User{}
	pbUser.Id = user.ID
	rsp.User = pbUser
	return nil
}

func (srv *UserServiceHandler) Get(ctx context.Context, req *pb.UserRequest, rsp *pb.UserResponse) error {
	rsp.User.Id = 1
	rsp.User.Name = "test"
	return nil
}

func (srv *UserServiceHandler) GetAll(ctx context.Context, req *pb.UserRequest, rsp *pb.UserResponse) error {
	user := &pb.User{
		Id:       0,
		Name:     "",
		Email:    "",
		Phone:    "",
		RealName: "",
		Avatar:   "",
		Status:   0,
		CreateAt: "",
		UpdateAt: "",
	}
	rsp.User = user

	return nil
}

func (srv *UserServiceHandler) Create(ctx context.Context, req *pb.UserRequest, rsp *pb.UserResponse) error {
	rsp.User.Id = 1
	rsp.User.Name = "test"
	return nil
}

func (srv *UserServiceHandler) Delete(ctx context.Context, req *pb.UserRequest, rsp *pb.UserResponse) error {
	rsp.User.Id = 1
	rsp.User.Name = "test"
	return nil
}

func (srv *UserServiceHandler) Update(ctx context.Context, req *pb.UserRequest, rsp *pb.UserResponse) error {
	rsp.User.Id = 1
	rsp.User.Name = "test"
	return nil
}
