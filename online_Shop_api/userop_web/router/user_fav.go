package router

import (
	"github.com/gin-gonic/gin"
	"online_Shop_api/userop_web/api/user_fav"
	"online_Shop_api/userop_web/middleware"
)

func InitUserFavRouter(Router *gin.RouterGroup) {
	UserFavRouter := Router.Group("userfav")
	{
		UserFavRouter.GET("", middleware.JwtToken(), user_fav.List)
		UserFavRouter.POST("", middleware.JwtToken(), user_fav.New)
		UserFavRouter.DELETE("/:id", middleware.JwtToken(), user_fav.Delete)
		UserFavRouter.GET("/:id", middleware.JwtToken(), user_fav.Detail)
	}
}
