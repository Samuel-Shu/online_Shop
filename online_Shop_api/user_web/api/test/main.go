package main

import (
	"online_Shop_api/user_web/api"
	"online_Shop_api/user_web/initialize"
)

func main() {
	initialize.InitConfig()
	api.SendEmail("2385970514@qq.com")
}
