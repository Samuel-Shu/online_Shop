package handler

import (
	"context"
	_ "crypto/md5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"online_Shop/inventory_srv/model"
	"online_Shop/inventory_srv/proto"
	"online_Shop/user_srv/global"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (i *InventoryServer) SetInv(c context.Context, req *proto.GoodsInvInfo) (*proto.MyEmpty, error) {
	//设置库存，
	var inv model.Inventory
	global.DB.First(&inv, req.GoodsId)
	inv.Goods = req.GoodsId
	inv.Stocks = uint32(req.Num)

	global.DB.Save(&inv)

	return &proto.MyEmpty{}, nil
}

func (i *InventoryServer) InvDetail(c context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.First(&inv, req.GoodsId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num: int32(inv.Stocks),
	}, nil
}

func (i *InventoryServer) Sell(c context.Context, req *proto.SellInfo) (*proto.MyEmpty, error) {
	//1、此处需要考虑到事务
	//如果购物车存在三个商品，数量分别是[1:15, 2:6, 3:7]
	//此时如果第一个商品数量满足并且扣减成功，但是第二个商品数量不足无法扣减
	//那么此处应当是购买失败的，第一个商品的数量也不应该扣减，应当回滚到扣减之前
	//2、并发情况下可能会出现超卖问题
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		//判断能否扣减库存
		if result := global.DB.First(&inv, goodInfo.GoodsId); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		if inv.Stocks < uint32(goodInfo.Num) {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减库存， 并发时会出现数据不一致的问题， 需要使用分布式锁解决此类问题
		inv.Stocks -= uint32(goodInfo.Num)
		tx.Save(&inv)
	}
	tx.Commit()
	return &proto.MyEmpty{}, nil
}

func (i *InventoryServer) Reback(c context.Context, req *proto.SellInfo) (*proto.MyEmpty, error) {
	//库存归还：1、订单超时归还；2、订单创建失败，归还之前扣减的库存；3、手动归还
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		//判断能否扣减库存
		if result := global.DB.First(&inv, goodInfo.GoodsId); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//添加库存， 并发时会出现数据不一致的问题， 需要使用分布式锁解决此类问题
		inv.Stocks += uint32(goodInfo.Num)
		tx.Save(&inv)
	}
	tx.Commit()
	return &proto.MyEmpty{}, nil
}
