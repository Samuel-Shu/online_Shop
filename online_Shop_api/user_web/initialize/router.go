package initialize

import (
	"github.com/gin-gonic/gin"
	"online_Shop_api/user_web/router"
)

func Router() *gin.Engine {
	Router := gin.Default()

	ApiGroup := Router.Group("/v1")
	router.InitUserRouter(ApiGroup)

	return Router
}
