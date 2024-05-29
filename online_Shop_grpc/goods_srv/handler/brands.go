package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"online_Shop/goods_srv/model"
	"online_Shop/goods_srv/proto"
	"online_Shop/user_srv/global"
)



func (g *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error){
	brandListResponse := proto.BrandListResponse{}
	var brands []model.Brands

	r := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if r.Error != nil {
		return nil, r.Error
	}

	var total int64
	global.DB.Model(&brands).Count(&total)

	brandListResponse.Total = int32(total)
	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id: int32(brand.ID),
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Data = brandResponses
	return &brandListResponse, nil
}

func (g *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error){
	//先判断是否已经存在该品牌
	if r := global.DB.First(&model.Brands{}); r.RowsAffected != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}
	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(brand)

	return &proto.BrandInfoResponse{Id: int32(brand.ID)}, nil
}

func (g *GoodsServer) DeleteBrand(c context.Context, req *proto.BrandRequest) (*proto.MyEmpty, error){
	if res := global.DB.Delete(&model.Brands{}, req.Id); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &proto.MyEmpty{}, nil
}

func (g *GoodsServer) UpdateBrand(c context.Context, req *proto.BrandRequest) (*proto.MyEmpty, error){
	//先判断是否已经存在该品牌
	brands := model.Brands{}
	if r := global.DB.First(&brands); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}

	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	global.DB.Save(&brands)
	return &proto.MyEmpty{}, nil
}
