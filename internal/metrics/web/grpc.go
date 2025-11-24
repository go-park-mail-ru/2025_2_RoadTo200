package web

import (
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
)

func NewGrpcServerInterceptor() (grpc.UnaryServerInterceptor, error) {
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
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	return traceutils.OpenTracingServerInterceptor(tracer), nil
}

func NewGrpcClientInterceptor() (grpc.UnaryClientInterceptor, error) {
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
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	return traceutils.OpenTracingClientInterceptor(tracer), nil
}
