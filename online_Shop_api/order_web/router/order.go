package router

import (
	"github.com/gin-gonic/gin"
	"online_Shop_api/order_web/api/order"
	"online_Shop_api/order_web/api/pay"
	"online_Shop_api/order_web/middleware"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("orders").Use(middleware.JwtToken())
	{
		OrderRouter.GET("", order.List)       //订单列表接口
		OrderRouter.POST("", order.New)       //新建订单接口
		OrderRouter.GET("/:id", order.Detail) //获取订单详情接口
	}

	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("alipay/notify", pay.Notify)
	}
}
