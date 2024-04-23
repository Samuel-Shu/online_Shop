package initialize

import (
	"github.com/gin-gonic/gin"
	"online_Shop_api/user_web/middleware"
	"online_Shop_api/user_web/router"
)

func Router() *gin.Engine {
	Router := gin.Default()
	//配置允许跨域请求
	Router.Use(middleware.Cors())
	ApiGroup := Router.Group("/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return Router
}
