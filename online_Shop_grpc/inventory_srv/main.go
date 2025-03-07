package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"online_Shop/inventory_srv/global"
	"online_Shop/inventory_srv/handler"
	"online_Shop/inventory_srv/initialize"
	"online_Shop/inventory_srv/proto"
	"online_Shop/inventory_srv/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "IP地址")
	PORT := flag.Int("port", 0, "port端口")
	flag.Parse()

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDb()
	initialize.InitRedisLock()

	zap.S().Info("ip地址：", *IP)
	if *PORT == 0 {
		*PORT, _ = utils.GetFreePort()
	}

	zap.S().Info("port端口：", *PORT)

	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *PORT))
	if err != nil {
		panic("端口监听失败" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Errorf("连接consul失败：%s", err)
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", "192.168.220.1", *PORT),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serverId := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serverId
	registration.Port = *PORT
	registration.Tags = []string{"onlineShop", "Samuel-Shu", "inventory", "srv"}
	registration.Address = "192.168.220.1"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic("grpc启动失败" + err.Error())
		}
	}()

	// 监听库存归还topic
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.220.128:9876"}),
		consumer.WithGroupName("onlineshop-inventory"),
	)
	if err := c.Subscribe("order_reback", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Println("获取到值：", msgs[i].Message.Body)
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		panic("订阅消息失败")
	}

	err = c.Start()
	if err != nil {
		return
	}

	time.Sleep(time.Second)
	err = c.Shutdown()
	if err != nil {
		return
	}

	//接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := client.Agent().ServiceDeregister(serverId); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
