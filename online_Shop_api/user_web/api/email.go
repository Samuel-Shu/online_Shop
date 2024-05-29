package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"net/http"
	"online_Shop_api/user_web/forms"
	"online_Shop_api/user_web/global"
	"online_Shop_api/user_web/utils"
	"time"
)

func SendEmail(c *gin.Context) {
	SendEmailForm := forms.SendEmail{}
	if err := c.ShouldBind(&SendEmailForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	captcha := utils.RandomCode()
	m := gomail.NewMessage()
	m.SetHeader("From", global.ServerConfig.EmailConfigInfo.FromEmail)
	m.SetHeader("To", SendEmailForm.Email)
	m.SetHeader("Subject", "登录验证码：")
	m.SetBody("test/html", "<p>【online_shop】的验证码是</p> <b>"+captcha+"</b>")
	//m.Attach("/home/Alex/lolcat.jpg")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisConfigInfo.IP, global.ServerConfig.RedisConfigInfo.Port),
	})

	d := gomail.NewDialer("smtp.qq.com", 587, global.ServerConfig.EmailConfigInfo.FromEmail, global.ServerConfig.EmailConfigInfo.SecretKey)

	//开启一个协程同步到redis中，尽可能保证验证码时效性
	go func(client *redis.Client) {
		err := rdb.Set(context.Background(), SendEmailForm.Email, captcha, time.Duration(global.ServerConfig.RedisConfigInfo.Expire)*time.Second).Err()
		if err != nil {

			panic(err)
		}
	}(rdb)
	// Send the emailFun to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	zap.S().Infof("【SendEmail】验证码生成成功：%s", captcha)
	c.JSON(http.StatusOK, gin.H{
		"msg": "发送验证码成功",
	})
}
