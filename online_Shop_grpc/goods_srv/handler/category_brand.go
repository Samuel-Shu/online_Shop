package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"online_Shop/goods_srv/model"
	"online_Shop/goods_srv/proto"
	"online_Shop/user_srv/global"
)

func (g *GoodsServer) CategoryBrandList(c context.Context, req *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	var categoryBrands []model.GoodsCategoryBrand
	CategoryBrandListResponse := proto.CategoryBrandListResponse{}

	var total int64
	global.DB.Model(&model.GoodsCategoryBrand{}).Count(&total)
	CategoryBrandListResponse.Total = int32(total)

	global.DB.Preload("Category").Preload("Brands ").Scopes(Paginate(int(req.Pages), int(req.PagePerNum))).Find(&categoryBrands)

	var categoryResponses []*proto.CategoryBrandResponse
	for _, categoryBrand := range categoryBrands {
		categoryResponses = append(categoryResponses, &proto.CategoryBrandResponse{
			Category: &proto.CategoryInfoResponse{
				Id:             int32(categoryBrand.Category.ID),
				Name:           categoryBrand.Category.Name,
				Level:          categoryBrand.Category.Level,
				IsTab:          categoryBrand.Category.IsTab,
				ParentCategory: int32(categoryBrand.Category.ParentCategoryID),
			},
			Brand: &proto.BrandInfoResponse{
				Id:   int32(categoryBrand.Brands.ID),
				Name: categoryBrand.Brands.Name,
				Logo: categoryBrand.Brands.Logo,
			},
		})
	}

	CategoryBrandListResponse.Data = categoryResponses
	return &CategoryBrandListResponse, nil
}

// 通过category获取brands

func (g *GoodsServer) GetCategoryBrandList(c context.Context, req *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	brandListResponse := proto.BrandListResponse{}

	var category model.Category
	if r := global.DB.Find(&category, req.Id).First(&category); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var categoryBrands []model.GoodsCategoryBrand
	if r := global.DB.Where(&model.GoodsCategoryBrand{CategoryID: category.ID}).Find(&categoryBrands); r.RowsAffected > 0 {
		brandListResponse.Total = int32(r.RowsAffected)
	}

	var brandInfoResponses []*proto.BrandInfoResponse
	for _, categoryBrand := range categoryBrands {
		brandInfoResponses = append(brandInfoResponses, &proto.BrandInfoResponse{
			Id:   int32(categoryBrand.Brands.ID),
			Name: categoryBrand.Brands.Name,
			Logo: categoryBrand.Brands.Logo,
		})
	}

	brandListResponse.Data = brandInfoResponses
	return &brandListResponse, nil
}
func (g *GoodsServer) CreateCategoryBrand(c context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	var category model.Category
	if r := global.DB.First(&category, req.CategoryId); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if r := global.DB.First(&brand, req.BrandId); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	categoryBrand := model.GoodsCategoryBrand{
		CategoryID: uint(req.CategoryId),
		BrandsID:   uint(req.BrandId),
	}
	global.DB.Save(&categoryBrand)
	return &proto.CategoryBrandResponse{Id: int32(categoryBrand.ID)}, nil
}

func (g *GoodsServer) DeleteCategoryBrand(c context.Context, req *proto.CategoryBrandRequest) (*proto.MyEmpty, error) {
	if res := global.DB.Delete(&model.GoodsCategoryBrand{}, req.Id); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌分类不存在")
	}
	return &proto.MyEmpty{}, nil
}

func (g *GoodsServer) UpdateCategoryBrand(c context.Context, req *proto.CategoryBrandRequest) (*proto.MyEmpty, error) {
	var categoryBrand model.GoodsCategoryBrand
	if r := global.DB.First(&categoryBrand, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌分类不存在")
	}

	var category model.Category
	if r := global.DB.First(&category, req.CategoryId); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if r := global.DB.First(&brand, req.BrandId); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	categoryBrand.CategoryID = uint(req.CategoryId)
	categoryBrand.BrandsID = uint(req.BrandId)
	global.DB.Save(&categoryBrand)

	return &proto.MyEmpty{}, nil
}
