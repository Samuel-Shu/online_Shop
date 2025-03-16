package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"online_Shop_api/llms_web/global"
	"online_Shop_api/llms_web/proto"
	"online_Shop_api/llms_web/utils/otgrpc"
)

func InitSrvConn() {
	fmt.Println(global.ServerConfig.ConsulInfo)
	llmsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.LlmsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"localBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("【InitSrvConn】连接【商品推荐服务失败】")
	}

	//生成grpc的client并调用接口
	global.LlmsSrvInfo = proto.NewLlmsClient(llmsConn)

	fmt.Println(global.LlmsSrvInfo)

}
