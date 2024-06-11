package model

import (
	"gorm.io/gorm"
	"time"
)

type ShoppingCart struct {
	gorm.Model
	User    int32 `gorm:"type:int;index"` //在购物车列表中需要查询当前用户的购物车记录
	Goods   int32 `gorm:"type:int;index"`
	Nums    int32 `gorm:"type:int"` //表示该购物车中该商品的数量
	Checked bool  //该商品是否选中（发生在购物车预支付时，选中购物车的某些商品进行支付，由checked字段判定该商品是否需要支付）
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

// OrderInfo 订单表
type OrderInfo struct {
	gorm.Model

	User    int32  `gorm:"type:int;index"`
	OrderSn string `gorm:"type:varchar(30);index"`                             //订单号，由平台自动生成
	PayType string `gorm:"type:varchar(20) comment 'alipay(支付宝), wechat(微信)'"` //选择支付类型

	Status     string    `gorm:"type:varchar(20) comment 'PAYING(待支付), TRADE_SUCCESS(成功), TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	TradeNO    string    `gorm:"type:varchar(100) comment '交易号'"` //交易号就是支付宝的订单号
	OrderMount float32   //订单金额
	PayTime    time.Time //用户支付时间

	Address      string `gorm:"type:varchar(100)"` //收货地址
	SignerName   string `gorm:"type:varchar(20)"`  //收件人姓名
	SignerMobile string `gorm:"type:varchar(11)"`  //收件人电话
	Post         string `gorm:"type:varchar(20)"`  //留言信息
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

type OrderGoods struct {
	gorm.Model

	Order      int32  `gorm:"type:int;index"`
	GoodsName  string `gorm:"type:varchar(100);index"`
	Goods      int32  `gorm:"type:int"`
	GoodsImage string `gorm:"type:varchar(200)"`
	GoodsPrice float32
	Nums       int32 `gorm:"type:int"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}
