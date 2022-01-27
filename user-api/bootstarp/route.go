package bootstarp

import (
	"github.com/869413421/pg-service/user-api/routes"
	"github.com/gin-gonic/gin"
	"sync"
)

var Router *gin.Engine
var once sync.Once

func SetupRoute() *gin.Engine {
	once.Do(func() {
		Router = gin.New()
		routes.RegisterWebRoutes(Router)
	})

	return Router
}
