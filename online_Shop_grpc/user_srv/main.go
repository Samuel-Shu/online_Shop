package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"online_Shop/user_srv/handler"
	"online_Shop/user_srv/proto"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "IP地址")
	PORT := flag.Int("port", 8081, "port端口")
	flag.Parse()
	fmt.Println("ip地址：", *IP)
	fmt.Println("port端口：", *PORT)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *PORT))
	if err != nil {
		panic("端口监听失败" + err.Error())
	}
	err = server.Serve(listen)
	if err != nil {
		panic("grpc启动失败" + err.Error())
	}
}
