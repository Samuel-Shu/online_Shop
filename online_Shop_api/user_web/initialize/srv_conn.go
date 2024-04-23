package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"online_Shop_api/user_web/global"
	"online_Shop_api/user_web/proto"
)

func InitSrvConn() {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	var (
		userSrvHost string
		userSrvPort int
	)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.Name))
	if err != nil {
		panic(err)
	}
	for _, v := range data {
		userSrvHost = v.Address
		userSrvPort = v.Port
		break
	}

	if userSrvHost == "" {
		zap.S().Fatal("【InitSrvConn】连接【用户服务失败】")
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[Register] 连接【用户服务】失败",
			"msg", err.Error(),
		)
	}
	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}