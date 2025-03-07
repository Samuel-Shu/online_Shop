package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"time"
)

func main() {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.220.128:9876"}),
		consumer.WithGroupName("online"),
	)
	if err != nil {
		panic("连接失败")
	}

	if err := c.Subscribe("onlineshop", consumer.MessageSelector{}, func(ctx context.Context,
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
}
