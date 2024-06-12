package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"online_Shop_api/order_web/global"
	"strings"
)

type JWT struct {
	JwtKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey),
	}
}

type MyClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}

// 定义错误
var (
	TokenExpired     error = errors.New("token已过期,请重新登录")
	TokenNotValidYet error = errors.New("token无效,请重新登录")
	TokenMalformed   error = errors.New("token不正确,请重新登录")
	TokenInvalid     error = errors.New("这不是一个token,请重新登录")
)

// CreateToken 生成token
func (j *JWT) CreateToken(claims MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey)
}

// ParserToken 解析token
func (j *JWT) ParserToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.JwtKey, nil
		})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if token != nil {
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	}

	return nil, TokenInvalid
}

// JwtToken jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "token不存在",
			})
			c.Abort()
			return
		}

		checkToken := strings.Split(tokenHeader, " ")
		if len(checkToken) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "1、token错误",
			})
			c.Abort()
			return
		}

		//if len(checkToken) != 2 || checkToken[0] != "Bearer" {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"msg": "2、token错误",
		//	})
		//	c.Abort()
		//	return
		//}

		j := NewJWT()
		// 解析token
		claims, err := j.ParserToken(checkToken[0]) //若需要考虑是否有Bearer字段，则此处应该是checkToken[1]
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "token授权已过期,请重新登录",
				})
				c.Abort()
				return
			}
			// 其他错误
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "3、token错误",
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}
