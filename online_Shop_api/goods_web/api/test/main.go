package main

import (
	"online_Shop_api/goods_web/api"
	"online_Shop_api/goods_web/initialize"
)

func main() {
	initialize.InitConfig()
	api.SendEmail("2385970514@qq.com")
}
