package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.220.128:9876"}), producer.WithRetry(2))
	if err != nil {
		panic(err)
	}
	if err := p.Start(); err != nil {
		panic("启动producer失败")
	}

	defer func(p rocketmq.Producer) {
		err := p.Shutdown()
		if err != nil {
			panic(err)
		}
	}(p)

	sync, err := p.SendSync(context.Background(), &primitive.Message{Topic: "onlineshop", Body: []byte("this a is onlineshop")})
	if err != nil {
		panic("sendSync失败")
	} else {
		fmt.Println(sync.String())
	}

}
