package handler

import (
	"context"
	pb "github.com/869413421/pg-service/user/proto"
)

type UserServiceHandler struct {
}

func NewUserServiceHandler() *UserServiceHandler {
	return &UserServiceHandler{}
}

func (srv *UserServiceHandler) Get(ctx context.Context, req *pb.UserRequest, rsp *pb.UserResponse) error {
	rsp.User.Id = 1
	rsp.User.Name = "test"
	return nil
}

func (srv *UserServiceHandler) GetAll(ctx context.Context, req *pb.UserRequest, rsp *pb.UserResponse) error {
	rsp.User.Id = 1
	rsp.User.Name = "test"
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
