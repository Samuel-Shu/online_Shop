package global

import (
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	_ "gorm.io/gorm/schema"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:sx221410@tcp(192.168.137.131:3306)/onlineShop_user_srv?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                                                 // string 类型字段的默认长度
		DontSupportRenameIndex:    true,                                                                                                // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                                // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                                               // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)

	}
}
