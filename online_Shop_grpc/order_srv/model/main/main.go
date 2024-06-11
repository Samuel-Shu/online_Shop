package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	_ "online_Shop/order_srv/global"
	"online_Shop/order_srv/model"
)

// 用于第一次部署项目时创建数据库表使用
// 使用gorm的autoMigrate函数进行表迁移
func main() {
	var err error
	DSN := "root:sx221410@tcp(192.168.137.134:3306)/onlineShop_order_srv?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       DSN,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		panic(err)

	}
	err = db.AutoMigrate(
		&model.OrderInfo{},
		&model.OrderGoods{},
		&model.ShoppingCart{},
	)
	if err != nil {
		panic(error.Error(err))
	}
}
