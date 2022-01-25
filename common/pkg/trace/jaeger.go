package trace

import (
	"github.com/go-acme/lego/v3/platform/config/env"
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
