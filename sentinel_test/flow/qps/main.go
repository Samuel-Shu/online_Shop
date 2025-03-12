package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"
	"math/rand"
	"time"
)

func main() {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("初始化sentinel异常:%v", err)
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "flow-qps-resource",
			TokenCalculateStrategy: flow.WarmUp,
			ControlBehavior:        flow.Reject,
			Threshold:              1000,
			StatIntervalInMs:       1000,
			WarmUpPeriodSec:        30,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for {
				e, b := sentinel.Entry("flow-qps-resource", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					// Blocked. We could get the block reason from the BlockError.
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					fmt.Println("限流")
				} else {
					// Passed, wrap the logic here.
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					fmt.Println("检查通过")

					// Be sure the entry is exited finally.
					e.Exit()
				}

			}
		}()
	}
	<-ch
}
