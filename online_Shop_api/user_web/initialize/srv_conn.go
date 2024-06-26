package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"online_Shop_api/user_web/global"
	"online_Shop_api/user_web/proto"
)

func InitSrvConn() {
	fmt.Println(global.ServerConfig.ConsulInfo)
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"localBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("【InitSrvConn】连接【用户服务失败】")
	}
	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
	fmt.Println(global.UserSrvClient)
}
