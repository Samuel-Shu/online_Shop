package global

import (
	"github.com/go-redsync/redsync/v4"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "gorm.io/gorm/schema"
	"online_Shop/order_srv/config"
	"online_Shop/order_srv/proto"
)

var (
	DB                 *gorm.DB
	ServerConfig       = &config.ServerConfig{}
	NacosConfig        = &config.NacosConfig{}
	RDBLock            *redsync.Redsync
	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)
