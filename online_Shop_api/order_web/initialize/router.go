package initialize

import (
	"github.com/gin-gonic/gin"
	"online_Shop_api/order_web/middleware"
	"online_Shop_api/order_web/router"
)

func Router() *gin.Engine {
	Router := gin.Default()
	//配置允许跨域请求
	Router.Use(middleware.Cors())
	ApiGroup := Router.Group("/o/v1")
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)

	return Router
}
