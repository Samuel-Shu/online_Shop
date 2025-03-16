package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online_Shop_api/inventory_web/middleware"
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
	_ = Router.Group("/inv/v1")

	return Router
}
