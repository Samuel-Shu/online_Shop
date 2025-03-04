package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"online_Shop_api/userop_web/global"
	"online_Shop_api/userop_web/proto"
)

func InitSrvConn() {
	fmt.Println(global.ServerConfig.ConsulInfo)
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"localBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("【InitSrvConn】连接【商品服务失败】")
	}

	userOpConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.UserOpSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"localBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("【InitSrvConn】连接【用户操作服务失败】")
	}

	//生成grpc的client并调用接口
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)
	global.UserFavSrvClient = proto.NewUserFavClient(userOpConn)
	global.AddressSrvClient = proto.NewAddressClient(userOpConn)
	global.MessageSrvClient = proto.NewMessageClient(userOpConn)
	fmt.Println(global.GoodsSrvClient)
	fmt.Println(global.UserFavSrvClient)
	fmt.Println(global.AddressSrvClient)
	fmt.Println(global.MessageSrvClient)
}
