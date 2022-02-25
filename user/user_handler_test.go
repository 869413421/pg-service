package main

import (
	"context"
	"github.com/869413421/pg-service/user/handler"
	"github.com/869413421/pg-service/user/pkg/model"
	"github.com/869413421/pg-service/user/pkg/repo/mocks"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/869413421/pg-service/user/service"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserServiceHandlerGet(t *testing.T) {
	ctx := context.Background()

	repo := &mocks.UserRepositoryInterface{}
	repo.On("GetByID", uint64(1)).Return(&model.User{Name: "Test", Email: "13528685024@163.com"}, nil)

	serviceHandler := &handler.UserServiceHandler{UserRepo: repo}
	serviceHandler.UserRepo = repo
	Convey("Testing UserServiceHandler Get Method", t, func() {
		req := &pb.GetRequest{Id: uint64(1)}
		rsp := &pb.UserResponse{}
		err := serviceHandler.Get(ctx, req, rsp)
		require.NoError(t, err)

		Convey("Then the response Name should be Test", func() {
			So(rsp.User.Name, ShouldEqual, "Test")
		})
	})
}

func TestUserServiceHandlerAuth(t *testing.T) {
	ctx := context.Background()

	repo := &mocks.UserRepositoryInterface{}
	repo.On("GetByEmail", "13528685024@163.com").Return(&model.User{Email: "13528685024@163.com", Password: "$2a$14$Rq.rwi2hAu4EDZ28qupN4uPZr7A3fPp.NcqWN7OUO9wbk5xAwA6vG"}, nil)
	tokenService := &service.TokenService{Repo: repo}
	serviceHandler := &handler.UserServiceHandler{UserRepo: repo, TokenService: tokenService}
	serviceHandler.UserRepo = repo
	Convey("Testing UserServiceHandler Auth Method", t, func() {
		req := &pb.AuthRequest{Email: "13528685024@163.com", Password: "123456"}
		rsp := &pb.TokenResponse{}
		err := serviceHandler.Auth(ctx, req, rsp)

		require.NoError(t, err)

		Convey("Then the response TokenLen should be 468", func() {
			So(len(rsp.Token), ShouldEqual, 468)
		})
	})
}
