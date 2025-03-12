package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("初始化sentinel异常:%v", err)
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              1000,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		zap.S().Fatalf("Unexpected error: %+v", err)
		return
	}
}
