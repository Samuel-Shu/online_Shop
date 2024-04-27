package global

import (
	ut "github.com/go-playground/universal-translator"
	"online_Shop_api/user_web/config"
	"online_Shop_api/user_web/proto"
)

var (
	Trans ut.Translator

	ServerConfig = &config.ServerConfig{}

	NacosConfig = &config.NacosConfig{}

	UserSrvClient proto.UserClient
)
