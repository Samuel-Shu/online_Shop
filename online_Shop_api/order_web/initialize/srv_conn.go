package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"online_Shop_api/order_web/global"
	"online_Shop_api/order_web/proto"
	"online_Shop_api/order_web/utils/otgrpc"
)

func InitSrvConn() {
	fmt.Println(global.ServerConfig.ConsulInfo)
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"localBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),

	)
	if err != nil {
		zap.S().Fatal("【InitSrvConn】连接【商品服务失败】")
	}

	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"localBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),

	)
	if err != nil {
		zap.S().Fatal("【InitSrvConn】连接【订单服务失败】")
	}

	InventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"localBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		
	)
	if err != nil {
		zap.S().Fatal("【InitSrvConn】连接【库存服务失败】")
	}

	//生成grpc的client并调用接口
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)
	global.OrderSrvClient = proto.NewOrderClient(orderConn)
	global.InventorySrvClient = proto.NewInventoryClient(InventoryConn)
	fmt.Println(global.GoodsSrvClient)
	fmt.Println(global.OrderSrvClient)
	fmt.Println(global.InventorySrvClient)
}
