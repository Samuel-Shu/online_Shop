package handler

import (
	"context"
	_ "crypto/md5"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"online_Shop/userop_srv/global"
	"online_Shop/userop_srv/model"
	"online_Shop/userop_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (i *InventoryServer) SetInv(c context.Context, req *proto.GoodsInvInfo) (*proto.MyEmptyWithInv, error) {
	//设置库存，
	var inv model.Inventory
	global.DB.First(&inv, req.GoodsId)
	inv.Goods = req.GoodsId
	inv.Stocks = uint32(req.Num)

	global.DB.Save(&inv)

	return &proto.MyEmptyWithInv{}, nil
}

func (i *InventoryServer) InvDetail(c context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.First(&inv, req.GoodsId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     int32(inv.Stocks),
	}, nil
}

func (i *InventoryServer) Sell(c context.Context, req *proto.SellInfo) (*proto.MyEmptyWithInv, error) {
	//1、此处需要考虑到事务
	//如果购物车存在三个商品，数量分别是[1:15, 2:6, 3:7]
	//此时如果第一个商品数量满足并且扣减成功，但是第二个商品数量不足无法扣减
	//那么此处应当是购买失败的，第一个商品的数量也不应该扣减，应当回滚到扣减之前
	//2、并发情况下可能会出现超卖问题
	//最终决定采用乐观锁机制：即采用version字段控制
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		//判断能否扣减库存
		mutex := global.RDBLock.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

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

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放分布式锁失败")
		}
	}
	tx.Commit()
	return &proto.MyEmptyWithInv{}, nil
}

func (i *InventoryServer) Reback(c context.Context, req *proto.SellInfo) (*proto.MyEmptyWithInv, error) {
	//库存归还：1、订单超时归还；2、订单创建失败，归还之前扣减的库存；3、手动归还
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory

		mutex := global.RDBLock.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		//判断能否扣减库存
		if result := global.DB.First(&inv, goodInfo.GoodsId); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//添加库存， 并发时会出现数据不一致的问题， 需要使用分布式锁解决此类问题
		inv.Stocks += uint32(goodInfo.Num)
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放分布式锁失败")
		}
	}
	tx.Commit()
	return &proto.MyEmptyWithInv{}, nil
}
