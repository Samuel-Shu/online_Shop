package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"online_Shop_api/user_web/forms"
	"online_Shop_api/user_web/global"
	"online_Shop_api/user_web/global/response"
	"online_Shop_api/user_web/middleware"
	"online_Shop_api/user_web/proto"
	"strconv"
	"strings"
	"time"
)

func RemoveTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转化为http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "内部错误",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(c *gin.Context) {
	//从表单获取数据
	pnStr := c.DefaultQuery("pn", "0")
	pn, _ := strconv.Atoi(pnStr)
	psizeStr := c.DefaultQuery("psize", "10")
	pSize, _ := strconv.Atoi(psizeStr)

	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pn),
		PSize: uint32(pSize),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询【用户列表】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		data := response.UserResponse{
			Id:       value.Id,
			Gender:   value.Gender,
			Email:    value.Email,
			NickName: value.NickName,
			Birthday: time.Unix(int64(value.Birthday), 0).Format("2006-01-02"),
			Role:     value.Role,
		}
		result = append(result, data)
	}

	c.JSON(http.StatusOK, result)
	zap.S().Debug("获取用户列表页")
}

func PasswordLogin(c *gin.Context) {
	//进行表单验证，确认传入数据有效
	PasswordLoginForm := forms.PasswordLoginForm{}
	if err := c.ShouldBind(&PasswordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	if !store.Verify(PasswordLoginForm.CaptchaId, PasswordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}


	if rsp, err := global.UserSrvClient.GetUserByEmail(context.Background(), &proto.EmailRequest{
		Email: PasswordLoginForm.Email,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"email": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"email": "登录失败",
				})
			}
			return
		}
	} else {
		//验证密码是否正确
		if passRsp, passErr := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          PasswordLoginForm.Password,
			EncryptedPassword: rsp.Password,
		}); passErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登陆失败",
			})
		} else {
			if passRsp.Success {
				//使用jwt生成token
				j := middleware.NewJWT()
				claims := middleware.MyClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               // 签名生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //设置30天过期
						Issuer:    "Samuel-Shu",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登陆失败",
				})
			}

		}
	}

}

func Register(c *gin.Context) {
	//进行表单验证，确认传入数据有效
	RegisterForm := forms.RegisterForm{}
	if err := c.ShouldBind(&RegisterForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisConfigInfo.IP, global.ServerConfig.RedisConfigInfo.Port),
	})

	result, err := rdb.Get(context.Background(), RegisterForm.Email).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	} else {
		if result != RegisterForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}
	}

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Email:    RegisterForm.Email,
		NickName: RegisterForm.Email,
		Password: RegisterForm.Password,
	})

	if err != nil {
 		zap.S().Errorf("[Register] 【新建用户】失败：%s", err)
		HandleGrpcErrorToHttp(err, c)
		return
	}

	//使用jwt生成token
	j := middleware.NewJWT()
	claims := middleware.MyClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //设置30天过期
			Issuer:    "Samuel-Shu",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
