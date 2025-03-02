package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"online_Shop/userop_srv/global"
	"online_Shop/userop_srv/model"
	"online_Shop/userop_srv/proto"
)

func (*UserOpServer) GetFavList(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	var rsp proto.UserFavListResponse
	var userFavs []model.UserFav
	var userFavList []*proto.UserFavResponse

	r := global.DB.Where(&model.UserFav{User: req.UserId, Goods: req.GoodsId}).Find(&userFavs)
	rsp.Total = int32(r.RowsAffected)

	for _, userFav := range userFavs {
		userFavList = append(userFavList, &proto.UserFavResponse{
			UserId:  userFav.User,
			GoodsId: userFav.Goods,
		})
	}
	rsp.Data = userFavList

	return &rsp, nil
}

func (*UserOpServer) AddUserFav(ctx context.Context, req *proto.UserFavRequest) (*proto.EmptyWithUserFav, error) {
	var userFav model.UserFav

	userFav.User = req.UserId
	userFav.Goods = req.GoodsId

	global.DB.Save(&userFav)

	return &proto.EmptyWithUserFav{}, nil
}

func (*UserOpServer) DeleteUserFav(ctx context.Context, req *proto.UserFavRequest) (*proto.EmptyWithUserFav, error) {
	var userFav model.UserFav
	// Unscoped 硬删除
	if r := global.DB.Unscoped().Where("goods = ? and user = ?", req.GoodsId, req.UserId).Delete(&userFav); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}

	return &proto.EmptyWithUserFav{}, nil
}

func (*UserOpServer) GetUserFavDetail(ctx context.Context, req *proto.UserFavRequest) (*proto.EmptyWithUserFav, error) {
	var userFav model.UserFav
	if r := global.DB.Where("goods = ? and user = ?", req.GoodsId, req.UserId).Find(&userFav); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}

	return &proto.EmptyWithUserFav{}, nil
}
