package handler

import (
	_ "crypto/md5"
	"online_Shop/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// 商品接口

//func (g *GoodsServer) GoodsList(context.Context, *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
//
//}
//
//// 用户提交订单有多个商品，需批量查询商品信息
//
//func (g *GoodsServer) BatchGetGoods(context.Context, *proto.BatchGoodsInfo) (*proto.GoodsListResponse, error) {
//
//}
//func (g *GoodsServer) CreateGoods(context.Context, *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
//
//}
//func (g *GoodsServer) DeleteGoods(context.Context, *proto.DeleteGoodsInfo) (*proto.MyEmpty, error) {
//
//}
//func (g *GoodsServer) UpdateGoods(context.Context, *proto.CreateGoodsInfo) (*proto.MyEmpty, error) {
//
//}
//func (g *GoodsServer) GetGoodsDetail(context.Context, *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
//
//}
