package global

import (
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "gorm.io/gorm/schema"
	"online_Shop/user_srv/config"
)

var (
	DB *gorm.DB
	ServerConfig = &config.ServerConfig{}
)


