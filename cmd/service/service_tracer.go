package service

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
)

// 使用jaeger
func newJaegerTracer() (tracer opentracing.Tracer, closer io.Closer, err error) {
	cfg := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  tracerJaegerType,  //固定采样
			Param: tracerJaegerParam, //1=全采样、0=不采样
		},
		Reporter: &jaegerConfig.ReporterConfig{
			//QueueSize:          200, // 缓冲区越大内存消耗越大,默认100
			LogSpans:           tracerJaegerLogSpans,
			LocalAgentHostPort: tracerJaegerHost,
		},
		ServiceName: fmt.Sprintf("%s.%s", serverName, namespace),
	}
	metricsFactory := prometheus.New()
	tracer, closer, err = cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger), jaegerConfig.Metrics(metricsFactory))
	if err != nil {
		return
	}
	opentracing.SetGlobalTracer(tracer)
	return
}
