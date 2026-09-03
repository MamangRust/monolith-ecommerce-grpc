package middleware

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// metadataCarrier adapts gRPC metadata.MD to the OTel propagation.TextMapCarrier
// interface so trace context can be injected/extracted through gRPC metadata.
type metadataCarrier struct {
	md metadata.MD
}

func (c *metadataCarrier) Get(key string) string {
	values := c.md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (c *metadataCarrier) Set(key, value string) {
	c.md.Set(key, value)
}

func (c *metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c.md))
	for key := range c.md {
		keys = append(keys, key)
	}
	return keys
}

// TraceUnaryServerInterceptor extracts trace context from incoming gRPC metadata
// and starts a server span for the invoked method, enabling end-to-end tracing
// from the API gateway down to the dependency services.
func TraceUnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		carrier := &metadataCarrier{md: md}
		ctx = propagator.Extract(ctx, carrier)

		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()

		resp, err := handler(ctx, req)

		span.SetAttributes(attribute.String("rpc.system", "grpc"))
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}

		return resp, err
	}
}

// TraceUnaryClientInterceptor injects the current trace context into outgoing
// gRPC metadata so downstream services can continue the same trace.
func TraceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	propagator := otel.GetTextMapPropagator()

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		carrier := &metadataCarrier{md: md}
		propagator.Inject(ctx, carrier)
		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
