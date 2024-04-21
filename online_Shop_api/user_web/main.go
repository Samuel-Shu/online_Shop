package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"online_Shop_api/user_web/global"
	"online_Shop_api/user_web/initialize"
	myvalidator "online_Shop_api/user_web/validator"
)

func main() {
	//初始化日志服务
	initialize.InitLogger()
	//初始化config配置
	initialize.InitConfig()
	//初始化路由
	router := initialize.Router()
	//初始化翻译器
	err:= initialize.InitTrans("zh")
	if err != nil {
		panic(err)
	}

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)

	if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
