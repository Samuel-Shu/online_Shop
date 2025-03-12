package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

func main() {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.220.128:6831",
		},
		ServiceName: "onlineShop",
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}

	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(closer)

	parentSpan := tracer.StartSpan("main")
	span1 := tracer.StartSpan("funcA", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 5000)

	span1.Finish()

	time.Sleep(time.Second * 2)
	span2 := tracer.StartSpan("funcB", opentracing.ChildOf(span1.Context()))
	time.Sleep(time.Second * 3)
	span2.Finish()

	parentSpan.Finish()
}
