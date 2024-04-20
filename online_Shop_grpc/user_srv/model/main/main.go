package main

import (
	"online_Shop/user_srv/global"
	_ "online_Shop/user_srv/global"
	"online_Shop/user_srv/model"
)

// 用于第一次部署项目时创建数据库表使用
// 使用gorm的autoMigrate函数进行表迁移
func main() {

	user := model.User{}
	err := global.DB.AutoMigrate(&user)
	if err != nil {
		panic(error.Error(err))
	}
}
