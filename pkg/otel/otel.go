package otel_pkg

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Config struct {
	ServiceName            string
	ServiceVersion         string
	Environment            string
	Endpoint               string
	Insecure               bool
	EnableRuntimeMetrics   bool
	RuntimeMetricsInterval time.Duration
	SamplingFraction       float64
}


type Telemetry struct {
	config         Config
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	loggerProvider *sdklog.LoggerProvider
}

func NewTelemetry(config Config) *Telemetry {
	if config.RuntimeMetricsInterval == 0 {
		config.RuntimeMetricsInterval = 15 * time.Second
	}

	return &Telemetry{
		config: config,
	}
}

func (t *Telemetry) Init(ctx context.Context) error {
	if err := t.InitTracer(ctx); err != nil {
		return fmt.Errorf("failed to initialize tracer: %w", err)
	}

	if err := t.InitMeter(ctx); err != nil {
		log.Printf("Warning: failed to initialize meter: %v", err)
	}

	if err := t.InitLogger(ctx); err != nil {
		log.Printf("Warning: failed to initialize logger: %v", err)
	}

	if t.config.EnableRuntimeMetrics {
		if err := t.InitRuntimeMetrics(); err != nil {
			log.Printf("Warning: failed to initialize runtime metrics: %v", err)
		} else {
			log.Printf("Runtime metrics initialized successfully (interval: %v)", t.config.RuntimeMetricsInterval)
		}
	}

	return nil
}

func (t *Telemetry) InitTracer(ctx context.Context) error {
	exporter, err := t.createTraceExporter(ctx)
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	res, err := t.createResource(ctx)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	samplingFraction := t.config.SamplingFraction
	if samplingFraction <= 0 {
		samplingFraction = 1.0
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplingFraction))),
	)


	t.tracerProvider = tp
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

func (t *Telemetry) InitMeter(ctx context.Context) error {
	exporter, err := t.createMetricExporter(ctx)
	if err != nil {
		return fmt.Errorf("failed to create metric exporter: %w", err)
	}

	res, err := t.createResource(ctx)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)

	t.meterProvider = mp
	otel.SetMeterProvider(mp)

	return nil
}

func (t *Telemetry) InitLogger(ctx context.Context) error {
	exporter, err := t.createLogExporter(ctx)
	if err != nil {
		return fmt.Errorf("failed to create log exporter: %w", err)
	}

	res, err := t.createResource(ctx)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)

	t.loggerProvider = lp

	return nil
}

func (t *Telemetry) InitRuntimeMetrics() error {
	err := runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to start runtime metrics: %w", err)
	}

	log.Printf("Go runtime metrics collection started")
	log.Printf("Collecting metrics: GC, Memory, Goroutines, CPU, etc.")

	return nil
}

func (t *Telemetry) GetLogger() *sdklog.LoggerProvider {
	return t.loggerProvider
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	var errs []error

	if t.tracerProvider != nil {
		if err := t.tracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("tracer provider shutdown: %w", err))
		}
	}

	if t.meterProvider != nil {
		if err := t.meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("meter provider shutdown: %w", err))
		}
	}

	if t.loggerProvider != nil {
		if err := t.loggerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("logger provider shutdown: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

func (t *Telemetry) createTraceExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(t.config.Endpoint),
	}

	if t.config.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	return otlptracegrpc.New(ctx, opts...)
}

func (t *Telemetry) createMetricExporter(ctx context.Context) (sdkmetric.Exporter, error) {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(t.config.Endpoint),
	}

	if t.config.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	return otlpmetricgrpc.New(ctx, opts...)
}

func (t *Telemetry) createLogExporter(ctx context.Context) (sdklog.Exporter, error) {
	opts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(t.config.Endpoint),
	}

	if t.config.Insecure {
		opts = append(opts, otlploggrpc.WithInsecure())
	}

	return otlploggrpc.New(ctx, opts...)
}

func (t *Telemetry) createResource(ctx context.Context) (*sdkresource.Resource, error) {
	return sdkresource.New(
		ctx,
		sdkresource.WithAttributes(
			semconv.ServiceNameKey.String(t.config.ServiceName),
			semconv.ServiceVersionKey.String(t.config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(t.config.Environment),
		),
	)
}
