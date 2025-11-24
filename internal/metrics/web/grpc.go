package web

import (
	"io"

	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"google.golang.org/grpc"
)

func NewGrpcServerInterceptor() (grpc.UnaryServerInterceptor, io.Closer, error) {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "session",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(prometheus.New()),
	)

	if err != nil {
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return traceutils.OpenTracingServerInterceptor(tracer), closer, nil
}

func NewGrpcClientInterceptor() (grpc.UnaryClientInterceptor, io.Closer, error) {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "client",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(prometheus.New()),
	)

	if err != nil {
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return traceutils.OpenTracingClientInterceptor(tracer), closer, nil
}
