package handler

import (
	"context"
	_ "crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"online_Shop/goods_srv/global"
	"online_Shop/goods_srv/model"
	"online_Shop/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func modelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              int32(goods.ID),
		CategoryId:      int32(goods.CategoryID),
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		FavNum:          goods.FavNum,
		SoldNum:         goods.SoldNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   int32(goods.Category.ID),
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   int32(goods.Brands.ID),
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

// 商品接口

func (g *GoodsServer) GoodsList(c context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//功能需求：关键词搜索、查询新品、查询热门商品、通过价格区间筛选、通过商品分类筛选
	goodsListResponse := &proto.GoodsListResponse{}

	q := elastic.NewBoolQuery()
	localDb := global.DB.Model(&model.Goods{})
	if req.KeyWords != "" {
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsHot {
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}

	if req.PriceMin > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		q = q.Filter(elastic.NewTermQuery("brand_id", req.Brand))
	}

	var subQuery string

	categoryIds := make([]interface{}, 0)
	if req.TopCategory > 0 {
		var category model.Category
		if r := global.DB.First(&category, req.TopCategory); r.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		if category.Level == 1 {
			subQuery = fmt.Sprintf("SELECT id FROM category WHERE parent_category_id IN (SELECT id FROM category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("SELECT id FROM category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("SELECT id FROM category WHERE  id=%d", req.TopCategory)
		}

		type Result struct {
			ID int32
		}
		var results []Result
		global.DB.Model(model.Category{}).Raw(subQuery).Scan(&results)
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}

		// 生成terms查询
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
	}

	if req.Pages <= 0 {
		req.Pages = 1
	}

	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	do, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil, err
	}

	goodsIds := make([]int32, 0)
	goodsListResponse.Total = int32(do.Hits.TotalHits.Value)
	for _, value := range do.Hits.Hits {
		goods := model.EsGoods{}
		_ = json.Unmarshal(value.Source, &goods)
		goodsIds = append(goodsIds, goods.ID)
	}

	// 查询id在某个数组中的值
	var goods []model.Goods
	r := localDb.Preload("Category").Preload("Brands").Find(&goods, goodsIds)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, good := range goods {
		goodsInfoResponse := modelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return goodsListResponse, nil
}

// BatchGetGoods 用户提交订单有多个商品，需批量查询商品信息
func (g *GoodsServer) BatchGetGoods(c context.Context, req *proto.BatchGoodsInfo) (*proto.GoodsListResponse, error) {
	var goods []model.Goods
	goodsListResponse := &proto.GoodsListResponse{}

	r := global.DB.Where(&goods, req.Id)
	for _, good := range goods {
		goodsInfoResponse := modelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(r.RowsAffected)
	return goodsListResponse, nil
}

func (g *GoodsServer) CreateGoods(c context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if r := global.DB.First(&category, req.CategoryId); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if r := global.DB.First(&brand, req.BrandId); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	global.DB.Save(&goods)

	return &proto.GoodsInfoResponse{
		Id: int32(goods.ID),
	}, nil
}

func (g *GoodsServer) DeleteGoods(c context.Context, req *proto.DeleteGoodsInfo) (*proto.MyEmpty, error) {
	if r := global.DB.Delete(&model.Goods{}, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品不存在")
	}

	return &proto.MyEmpty{}, nil
}

func (g *GoodsServer) UpdateGoods(c context.Context, req *proto.CreateGoodsInfo) (*proto.MyEmpty, error) {
	var goods model.Goods

	if r := global.DB.First(&goods, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if r := global.DB.First(&category, req.CategoryId); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if r := global.DB.First(&brand, req.BrandId); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	global.DB.Save(&goods)
	return &proto.MyEmpty{}, nil
}

func (g *GoodsServer) GetGoodsDetail(c context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if r := global.DB.First(&goods, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	goodsInfoResponse := modelToResponse(goods)
	return &goodsInfoResponse, nil
}
