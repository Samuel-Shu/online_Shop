package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"online_Shop_api/goods_web/api/goods"
	"online_Shop_api/user_web/middleware"
)

func InitGoodsRouter(router *gin.RouterGroup) {
	zap.S().Infof("配置用户相关URL")
	goodsGroup := router.Group("goods")

	{
		goodsGroup.GET("", goods.List) //商品列表
		goodsGroup.POST("", middleware.JwtToken(), middleware.IsAdminAuth(), goods.New)
		goodsGroup.GET("/:id", goods.Detail) //商品详情
		goodsGroup.DELETE("/:id", middleware.JwtToken(), middleware.IsAdminAuth(), goods.Delete) //删除商品
		goodsGroup.GET("/:id/stocks", goods.Stocks) //获取商品库存

		goodsGroup.PUT("/:id", middleware.JwtToken(), middleware.IsAdminAuth(), goods.Update)
		goodsGroup.PATCH("/:id", middleware.JwtToken(), middleware.IsAdminAuth(), goods.UpdateStatus)
	}

}
