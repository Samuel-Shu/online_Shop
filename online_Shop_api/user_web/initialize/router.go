package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online_Shop_api/user_web/middleware"
	"online_Shop_api/user_web/router"
)

func Router() *gin.Engine {
	Router := gin.Default()
	//配置允许跨域请求
	Router.Use(middleware.Cors())
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	ApiGroup := Router.Group("/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return Router
}
