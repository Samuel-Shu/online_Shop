package router

import (
	"github.com/gin-gonic/gin"
	"online_Shop_api/llms_web/api"
	"online_Shop_api/llms_web/middleware"
)

func InitLlmsRouter(Router *gin.RouterGroup) {
	LlmsRouter := Router.Group("llms").Use(middleware.JwtToken())
	{
		LlmsRouter.GET("/health", api.HealthCheck) //llms接口可用性检查
		LlmsRouter.POST("/chat", api.Chat)         // 对话接口
		LlmsRouter.POST("/upload", api.UploadFile) //文件上传
	}
}
