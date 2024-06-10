package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_Shop_api/goods_web/api/category"
)

func InitCategoryRouter(router *gin.RouterGroup) {
	zap.S().Infof("配置目录相关URL")
	categoryGroup := router.Group("categorys")

	{
		categoryGroup.GET("", category.List)
		categoryGroup.DELETE("/:id", category.Delete)
		categoryGroup.GET("/:id", category.Detail)
		categoryGroup.POST("", category.New)
		categoryGroup.PUT("/:id", category.Update)
	}

}
