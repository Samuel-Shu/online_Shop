package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_Shop_api/user_web/api"
)

func InitUserRouter(router *gin.RouterGroup) {
	zap.S().Infof("配置用户相关URL")
	userGroup := router.Group("user")

	userGroup.GET("/list", api.GetUserList)

}
