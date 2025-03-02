package model

import "gorm.io/gorm"

// LeavingMessages 留言字段
type LeavingMessages struct {
	gorm.Model

	User        int32  `gorm:"type:int;index"`
	MessageType int32  `gorm:"type:int comment '留言类型：1（留言），2（投诉），3（询问），4（售后），5（求购）'"`
	Subject     string `gorm:"type:varchar(100) comment '主题'"`
	Message     string
	File        string `gorm:"type:varchar(200)"`
}

func (LeavingMessages) TableName() string {
	return "leavingmessages"
}

// Address 用户下单地址
type Address struct {
	gorm.Model

	User         int32  `gorm:"type:int;index"`
	Province     string `gorm:"type:varchar(10)"`
	City         string `gorm:"type:varchar(10)"`
	District     string `gorm:"type:varchar(20)"`
	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"`
	SignerMobile string `gorm:"type:varchar(11)"`
}

// UserFav 用户收藏
type UserFav struct {
	gorm.Model

	// 用户和商品直接建立联合的唯一索引，一个用户对一件商品只能收藏一次
	User  int32 `gorm:"type:int;index:idx_user_goods,unique"`
	Goods int32 `gorm:"type:int;index:idx_user_goods,unique"`
}

func (UserFav) TableName() string {
	return "userfav"
}
