package gin2micro

import (
	"github.com/869413421/pg-service/common/pkg/api/http/controller"
	"github.com/869413421/pg-service/common/pkg/logger"
	"github.com/869413421/pg-service/common/pkg/container"
	pb "github.com/869413421/pg-service/user/proto/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var base = controller.NewBaseController()

func Jwt(content *gin.Context) {
	logger.Danger("this is test jwt")
	ctx, ok := ContextWithSpan(content)
	if ok == false {
		logger.Warning("user api user/auth get context err")
	}
	//1.获取token
	token := content.GetHeader("Authorization")
	if token != "" {
		tokenS := strings.Split(token, " ")
		token = tokenS[1]
	} else {
		token = content.Request.FormValue("token")
	}
	if token == "" {
		base.ResponseJson(content, http.StatusUnauthorized, "not found token", []string{})
		return
	}

	//2.验证token
	req := &pb.TokenRequest{Token: token}
	client := container.GetUserServiceClient()
	rsp, err := client.ValidateToken(ctx, req)
	if err != nil {
		base.ResponseJson(content, http.StatusUnauthorized, err.Error(), []string{})
		return
	}

	//3.通过验证
	if rsp.Valid != true {
		base.ResponseJson(content, http.StatusUnauthorized, "valid token error", []string{})
		return
	}

	content.Next()
}
