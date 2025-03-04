package model

import (
	"context"
	"gorm.io/gorm"
	"online_Shop/goods_srv/global"
	"strconv"
)

type Category struct {
	gorm.Model
	Name             string      `gorm:"type:varchar(20);not null comment '目录名称'" json:"name"`
	ParentCategoryID uint        `json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"` //foreignKey指明了外键，references指明了主键，那么在查询时可以根据这俩内容来自动填充数据
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"`
	IsTab            bool        `gorm:"default:false;not null comment '是否是tab标签'" json:"is_tab"`
}

type Brands struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null comment '品牌名称'"`
	Logo string `gorm:"type:varchar(200);default:'' comment '品牌logo'"`
}

type GoodsCategoryBrand struct {
	gorm.Model
	CategoryID uint `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category
	BrandsID   uint `gorm:"type:int;index:idx_category_brand,unique"`
	Brands     Brands
}

// TableName gorm生成表名称自定义
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type Banner struct {
	gorm.Model
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

type Goods struct {
	gorm.Model

	CategoryID uint `gorm:"type:int;not null"`
	Category   Category
	BrandsID   uint `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null comment '判断商品是否上架'"`
	ShipFree bool `gorm:"default:false;not null comment '判断商品是否免运费'"`
	IsNew    bool `gorm:"default:false;not null comment '判断商品是否是新品'"`
	IsHot    bool `gorm:"default:false;not null comment '判断商品是否属于热门商品'"`

	Name            string   `gorm:"type:varchar(50);not null comment '商品名称'"`
	GoodsSn         string   `gorm:"type:varchar(50);not null comment '商品编号'"`
	ClickNum        int32    `gorm:"type:int;default:0;not null comment '被点击次数'"`
	SoldNum         int32    `gorm:"type:int;default:0;not null comment '销量'"`
	FavNum          int32    `gorm:"type:int;default:0;not null comment '被收藏量'"`
	MarketPrice     float32  `gorm:"not null comment '商品实际价格'"`
	ShopPrice       float32  `gorm:"not null comment '售卖价格'"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null comment '商品简介'"`
	Images          GormList `gorm:"type:varchar(1000);not null"`
	DescImages      GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`
}

func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          int32(g.ID),
		CategoryID:  int32(g.CategoryID),
		BrandsID:    int32(g.BrandsID),
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}
