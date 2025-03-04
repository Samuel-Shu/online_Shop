package router

import (
	"github.com/gin-gonic/gin"
	"online_Shop_api/userop_web/api/message"
	"online_Shop_api/userop_web/middleware"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(middleware.JwtToken())
	{
		MessageRouter.GET("", message.List)
		MessageRouter.POST("", message.New)
	}
}
