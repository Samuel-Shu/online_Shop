package router

import (
	"github.com/gin-gonic/gin"
	"online_Shop_api/userop_web/api/address"
	"online_Shop_api/userop_web/middleware"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address")
	{
		AddressRouter.GET("", middleware.JwtToken(), address.List)
		AddressRouter.POST("", middleware.JwtToken(), address.New)
		AddressRouter.DELETE("/:id", middleware.JwtToken(), address.Delete)
		AddressRouter.PUT("/:id", middleware.JwtToken(), address.Update)
	}
}
