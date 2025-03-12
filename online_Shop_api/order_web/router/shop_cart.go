package router

import (
	"github.com/gin-gonic/gin"
	"online_Shop_api/order_web/api/shop_cart"
	"online_Shop_api/order_web/middleware"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	ShopCartRouter := Router.Group("shopcarts").Use(middleware.JwtToken()).Use(middleware.Trace())
	{
		ShopCartRouter.GET("", shop_cart.List)          //购物车列表
		ShopCartRouter.DELETE("/:id", shop_cart.Delete) //删除购物车条目
		ShopCartRouter.POST("", shop_cart.New)          //添加购物车条目
		ShopCartRouter.PATCH("/:id", shop_cart.Update)  //修改条目

	}
}
