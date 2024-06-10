package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_Shop_api/goods_web/api/brand"
)

func InitBrandRouter(router *gin.RouterGroup) {
	zap.S().Infof("配置品牌相关URL")
	BrandGroup := router.Group("brands")

	{
		BrandGroup.GET("", brand.List)
		BrandGroup.DELETE("/:id", brand.Delete)
		BrandGroup.POST("", brand.New)
		BrandGroup.PUT("/:id", brand.Update)
	}

	CategoryBrandRouter := router.Group("categorybrands")

	{
		CategoryBrandRouter.GET("", brand.CategoryBrandList)
		CategoryBrandRouter.DELETE("/:id", brand.DeleteCategoryBrand)
		CategoryBrandRouter.POST("", brand.NewCategoryBrand)
		CategoryBrandRouter.PUT("/:id", brand.UpdateCategoryBrand)
	}
}
