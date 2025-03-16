package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_Shop_api/inventory_web/api"
	"online_Shop_api/inventory_web/middleware"
)

func InitInventoryRouter(router *gin.RouterGroup) {
	zap.S().Infof("配置播报相关URL")
	InventoryGroup := router.Group("inventory")

	{
		InventoryGroup.GET("", api.List)
		InventoryGroup.DELETE("/:id", middleware.JwtToken(), middleware.IsAdminAuth(), api.Delete)
		InventoryGroup.POST("", middleware.JwtToken(), middleware.IsAdminAuth(), api.New)
		InventoryGroup.PUT("/:id", middleware.JwtToken(), middleware.IsAdminAuth(), api.Update)
	}

}
