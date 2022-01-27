package trace

import (
	"context"
	"github.com/869413421/pg-service/common/pkg/logger"
	"github.com/go-acme/lego/v3/platform/config/env"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

func NewTracer(serviceName string, addr string) (opentracing.Tracer, io.Closer, error) {
	collectionEndpoint := env.GetOrDefaultString("MICRO_TRACE_ENDPOINT", "http://jaeger:14268/api/traces")
	config := jaegerConfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			CollectorEndpoint:   collectionEndpoint,
		},
	}

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := config.NewTracer(
		jaegerConfig.Reporter(reporter),
	)

	return tracer, closer, err
}

// NewSpan 创建新的SPAN
func NewSpan(spanName string, ctx context.Context, req interface{}, rsp interface{}) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	var sp opentracing.Span
	wireContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	if err != nil {
		logger.Warning("create span err:", err)
	}
	// 创建新的 Span 并将其绑定到微服务上下文
	sp = opentracing.StartSpan(spanName, opentracing.ChildOf(wireContext))
	// 记录请求
	sp.SetTag("req", req)
	// 记录响应
	sp.SetTag("res", rsp)
	// 在函数返回 stop span 之前，统计函数执行时间
	sp.Finish()
}
