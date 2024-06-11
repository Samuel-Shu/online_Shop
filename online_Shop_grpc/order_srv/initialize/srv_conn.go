package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"online_Shop/order_srv/global"
	"online_Shop/order_srv/proto"
)

func InitSrvConn()  {
	//初始化第三方微服务的client
	fmt.Println(global.ServerConfig.ConsulInfo)
	GoodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"localBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("【InitSrvConn】连接【商品服务失败】")
	}

	InventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"localBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("【InitSrvConn】连接【库存服务失败】")
	}

	//生成grpc的client并调用接口
	global.GoodsSrvClient = proto.NewGoodsClient(GoodsConn)
	global.InventorySrvClient = proto.NewInventoryClient(InventoryConn)
	fmt.Println(global.GoodsSrvClient)
	fmt.Println(global.InventorySrvClient)
}
