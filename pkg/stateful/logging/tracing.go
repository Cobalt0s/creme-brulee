package logging

import (
	"context"
	"github.com/Cobalt0s/creme-brulee/pkg/stateful/config"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type serviceNameKey struct{}

func RegisterTracing(ctx context.Context, cfg *config.JaegerTraceConfig, exporter tracesdk.SpanExporter, serviceName string) (context.Context, func()) {
	ctx = context.WithValue(ctx, serviceNameKey{}, serviceName)
	log := ctxlogrus.Extract(ctx)
	if exporter == nil {
		log.Warn("tracer: down, no exporter is registered")
		return ctx, func() {}
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(newResource(serviceName)),
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(cfg.SamplingProbability)),
	)
	otel.SetTracerProvider(tp)

	log.Info("tracer: ready")
	return ctx, func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}

func newResource(serviceName string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			//TODO semconv.ServiceVersionKey.String("v0.1.0"),
		),
	)
	return r
}

func NewFileExporter(ctx context.Context) (tracesdk.SpanExporter, func()) {
	log := ctxlogrus.Extract(ctx)
	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		log.Fatal(err)
	}

	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(f),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
	if err != nil {
		f.Close()
		log.Fatal(err)
	}

	log.Info("tracer: file exporter is created")
	return exporter, func() {
		f.Close()
	}
}

func NewJaegerExporter(ctx context.Context, cfg *config.JaegerTraceConfig) tracesdk.SpanExporter {
	log := ctxlogrus.Extract(ctx)

	if cfg.Enabled {
		exp, err := jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(cfg.GetURL()),
			),
		)
		if err != nil {
			log.Fatalf("failed create Jaeger exporter %v", err)
		}
		log.Info("tracer: jaeger exporter is created")
		return exp
	}
	return nil
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	serviceName := ctx.Value(serviceNameKey{}).(string)
	return otel.Tracer(serviceName).Start(ctx, name)
}
