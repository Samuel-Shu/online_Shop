package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc  {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		currentUser := claims.(*MyClaims)

		if currentUser.ID != 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
