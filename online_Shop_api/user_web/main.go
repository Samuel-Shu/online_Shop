package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"online_Shop_api/user_web/global"
	"online_Shop_api/user_web/initialize"
	"online_Shop_api/user_web/utils"
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

	//初始化srv连接
	initialize.InitSrvConn()

	viper.AutomaticEnv()
	debug := viper.GetBool("ONLINE_SHOP_DEBUG")
	//debug表示本地测试环境
	//如果是debug，则端口号固定，方便测试，如果非debug模式，端口号随机获取可用
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}

	}

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("email", myvalidator.ValidateEmail)
		_ = v.RegisterTranslation("email", global.Trans, func(ut ut.Translator) error {
			return ut.Add("email", "{0} 非法的邮箱", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("email", fe.Field())
			return t
		})
	}

	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)

	if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
