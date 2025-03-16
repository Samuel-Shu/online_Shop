package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"online_Shop_api/inventory_web/global"
	"online_Shop_api/inventory_web/initialize"
	"online_Shop_api/inventory_web/utils"
	"online_Shop_api/inventory_web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//初始化日志服务
	initialize.InitLogger()
	//初始化config配置
	initialize.InitConfig()
	//初始化路由
	router := initialize.Router()
	//初始化翻译器
	err := initialize.InitTrans("zh")
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

	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("注册失败", err.Error())
	}

	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)

	go func() {
		if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败", err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = registerClient.DeRegister(serviceId)
	if err != nil {
		zap.S().Info("服务注销失败", err.Error())
	} else {
		zap.S().Info("服务注销成功")
	}
}
