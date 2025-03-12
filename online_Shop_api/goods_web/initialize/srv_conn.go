package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"online_Shop_api/goods_web/global"
	"online_Shop_api/goods_web/proto"
	"online_Shop_api/goods_web/utils/otgrpc"
)

func InitSrvConn() {
	fmt.Println(global.ServerConfig.ConsulInfo, global.ServerConfig.GoodsSrvInfo.Name)
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"localBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("【InitSrvConn】连接【用户服务失败】")
	}
	//生成grpc的client并调用接口
	global.GoodsSrvClient = proto.NewGoodsClient(userConn)
	fmt.Println(global.GoodsSrvClient)
}
