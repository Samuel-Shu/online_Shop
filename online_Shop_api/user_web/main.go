package main

import (
	"fmt"
	"go.uber.org/zap"
	"online_Shop_api/user_web/initialize"
)

func main() {

	initialize.InitLogger()
	router := initialize.Router()

	PORT := 7071

	zap.S().Debugf("启动服务器，端口：%d", PORT)

	if err := router.Run(fmt.Sprintf(":%d", PORT)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
