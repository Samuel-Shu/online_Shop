package global

import (
	ut "github.com/go-playground/universal-translator"
	"online_Shop_api/order_web/config"
	"online_Shop_api/order_web/proto"
)

var (
	Trans ut.Translator

	ServerConfig = &config.ServerConfig{}

	NacosConfig = &config.NacosConfig{}

	GoodsSrvClient proto.GoodsClient

	OrderSrvClient proto.OrderClient

	InventorySrvClient proto.InventoryClient
)
