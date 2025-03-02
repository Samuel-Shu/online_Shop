package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"online_Shop/goods_srv/global"
	"online_Shop/goods_srv/model"
	"online_Shop/goods_srv/proto"
)

// 轮播图

func (g *GoodsServer) BannerList(context.Context, *proto.MyEmpty) (*proto.BannerListResponse, error) {
	bannerListResponse := proto.BannerListResponse{}

	var banners []model.Banner
	res := global.DB.Find(&banners)
	bannerListResponse.Total = int32(res.RowsAffected)

	var bannerResponses []*proto.BannerResponse
	for _, banner := range banners {
		bannerResponses = append(bannerResponses, &proto.BannerResponse{
			Id:    int32(banner.ID),
			Index: banner.Index,
			Image: banner.Image,
			Url:   banner.Url,
		})
	}
	bannerListResponse.Data = bannerResponses
	return &bannerListResponse, nil
}

func (g *GoodsServer) CreateBanner(c context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := model.Banner{}

	banner.Image = req.Image
	banner.Index = req.Index
	banner.Url = req.Url

	global.DB.Save(&banner)

	return &proto.BannerResponse{Id: int32(banner.ID)}, nil
}

func (g *GoodsServer) DeleteBanner(c context.Context, req *proto.BannerRequest) (*proto.MyEmpty, error) {
	if res := global.DB.Delete(&model.Banner{}, req.Id); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	return &proto.MyEmpty{}, nil
}

func (g *GoodsServer) UpdateBanner(c context.Context, req *proto.BannerRequest) (*proto.MyEmpty, error) {
	//先判断是否已经存在该品牌
	banner := model.Banner{}
	if r := global.DB.First(&banner); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}

	global.DB.Save(&banner)
	return &proto.MyEmpty{}, nil
}
