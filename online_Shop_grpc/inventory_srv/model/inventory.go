package model

import (
	"gorm.io/gorm"
)

//type Stock struct {
//	gorm.Model
//	Name string
//	Address string
//}

type Inventory struct {
	gorm.Model
	Goods int32 `gorm:"type:int;index"`
	Stocks uint32 `gorm:"type:int"`
	//Stock Stock
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁

}

type InventoryHistory struct {

}
