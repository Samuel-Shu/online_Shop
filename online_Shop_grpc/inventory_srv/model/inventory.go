package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

//type Stock struct {
//	gorm.Model
//	Name string
//	Address string
//}

type GoodsDetail struct {
	Goods int32
	Num   int32
}
type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type Inventory struct {
	gorm.Model
	Goods  int32  `gorm:"type:int;index"`
	Stocks uint32 `gorm:"type:int"`
	//Stock Stock
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁

}

type DeliverHistory struct {
	Goods   int32  `gorm:"type:int;index"`
	Nums    int32  `gorm:"type:int"`
	OrderSn string `gorm:"type:varchar(200)"`
	Status  string `gorm:"type:varchar(200)"` // 1表示等待支付 2表示支付成功 3表示支付失败
}

type StockSellDetail struct {
	OrderSn string          `gorm:"type:varchar(200);index:idx_order_sn,unique;"`
	Status  int32           `gorm:"type:varchar(200)"` // 1表示已扣减，2表示已归还
	Detail  GoodsDetailList `gorm:"type:varchar(200)"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}
