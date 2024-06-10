package global

import (
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "gorm.io/gorm/schema"
	"online_Shop/inventory_srv/config"
)

var (
	DB *gorm.DB
	ServerConfig = &config.ServerConfig{}
	NacosConfig = &config.NacosConfig{}
)


