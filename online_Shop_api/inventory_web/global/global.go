package global

import (
	ut "github.com/go-playground/universal-translator"
	"online_Shop_api/inventory_web/config"
	"online_Shop_api/inventory_web/proto"
)

var (
	Trans ut.Translator

	ServerConfig = &config.ServerConfig{}

	NacosConfig = &config.NacosConfig{}

	InventorySrvClient proto.InventoryClient
)
