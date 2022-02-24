package test

import (
	"context"
	"github.com/869413421/pg-service/user/handler"
	"github.com/869413421/pg-service/user/pkg/model"
	"github.com/869413421/pg-service/user/pkg/repo/mocks"
	pb "github.com/869413421/pg-service/user/proto/user"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUserServiceHandlerGet(t *testing.T) {
	ctx := context.Background()

	repo :=&mocks.UserRepositoryInterface{}
	repo.On("Get",1).Return(&model.User{Name: "Test"})
	repo.On("GetByEmail","13528685024@163.com").Return(&model.User{Email: "13528685024@163.com"})

	serviceHandler := handler.NewUserServiceHandler()
	serviceHandler.UserRepo = repo
	Convey("Testing UserServiceHandler Get Method", t, func() {
		Convey("Using ID 1 Request", func() {
			req := &pb.GetRequest{Id: 1}
			rsp := &pb.UserResponse{}
			serviceHandler.Get(ctx, req, rsp)

			Convey("Then the response Name should be Test", func() {
				So(rsp.User.Name, ShouldEqual, "Test")
			})
		})
	})
}
