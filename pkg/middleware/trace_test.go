package middleware

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestMetadataCarrierGetSetKeys(t *testing.T) {
	md := metadata.New(nil)
	carrier := &metadataCarrier{md: md}

	carrier.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	if got := carrier.Get("traceparent"); got != "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01" {
		t.Fatalf("Get returned %q, want traceparent value", got)
	}

	if got := carrier.Get("missing"); got != "" {
		t.Fatalf("Get on missing key returned %q, want empty", got)
	}

	keys := carrier.Keys()
	if len(keys) != 1 || keys[0] != "traceparent" {
		t.Fatalf("Keys returned %v, want [traceparent]", keys)
	}
}

func TestTraceUnaryClientInterceptorInjectsTraceparent(t *testing.T) {
	// In production the global propagator is set during telemetry.Init; here we
	// set it explicitly so the interceptor produces a traceparent header.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create a real span so the TraceContext propagator injects a traceparent
	// header (a no-op span would not produce one).
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	ctx, span := tp.Tracer("test").Start(context.Background(), "test-span")
	defer span.End()

	interceptor := TraceUnaryClientInterceptor()

	var capturedMD metadata.MD
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		capturedMD, _ = metadata.FromOutgoingContext(ctx)
		return nil
	}

	err := interceptor(ctx, "/svc.Method", nil, nil, nil, invoker)
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}

	carrier := &metadataCarrier{md: capturedMD}
	if carrier.Get("traceparent") == "" {
		t.Fatal("traceparent was not injected into outgoing gRPC metadata")
	}

	// The injected header must round-trip through Extract.
	extracted := otel.GetTextMapPropagator().Extract(context.Background(), carrier)
	if extracted == nil {
		t.Fatal("extracted context is nil")
	}
}
