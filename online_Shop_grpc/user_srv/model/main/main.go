package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	_ "online_Shop/user_srv/global"
	"online_Shop/user_srv/model"
)

// 用于第一次部署项目时创建数据库表使用
// 使用gorm的autoMigrate函数进行表迁移
func main() {

	user := model.User{}

	DSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		"root", "sx221410", "192.168.220.128", 3306, "onlineshop_user_srv")

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
	//err := global.DB.AutoMigrate(&user)
	//if err != nil {
	//	panic(error.Error(err))
	//}
	err = db.AutoMigrate(&user)
	if err != nil {
		panic(error.Error(err))
	}

}
