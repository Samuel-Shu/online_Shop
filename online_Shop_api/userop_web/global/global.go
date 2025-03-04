package global

import (
	ut "github.com/go-playground/universal-translator"
	"online_Shop_api/userop_web/config"
	"online_Shop_api/userop_web/proto"
)

var (
	Trans ut.Translator

	ServerConfig = &config.ServerConfig{}

	NacosConfig = &config.NacosConfig{}

	GoodsSrvClient proto.GoodsClient

	MessageSrvClient proto.MessageClient

	AddressSrvClient proto.AddressClient

	UserFavSrvClient proto.UserFavClient
)
