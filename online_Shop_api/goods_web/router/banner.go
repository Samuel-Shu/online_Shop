package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_Shop_api/goods_web/api/banner"
	"online_Shop_api/goods_web/middleware"
)

	func InitBannerRouter(router *gin.RouterGroup) {
		zap.S().Infof("配置播报相关URL")
		BannerGroup := router.Group("banners")

		{
			BannerGroup.GET("", banner.List)
			BannerGroup.DELETE("/:id", middleware.JwtToken(), middleware.IsAdminAuth(), banner.Delete)
			BannerGroup.POST("", middleware.JwtToken(), middleware.IsAdminAuth(), banner.New)
			BannerGroup.PUT("/:id",middleware.JwtToken(), middleware.IsAdminAuth(), banner.Update)
		}

	}
