package router

import (
	"github.com/gin-gonic/gin"
	"online_Shop_api/user_web/api"
)

func InitBaseRouter(Router *gin.RouterGroup)  {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
		BaseRouter.POST("send_email", api.SendEmail)
	}
}
