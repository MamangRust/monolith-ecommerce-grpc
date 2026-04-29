package middleware

import (
	"context"

	"github.com/grafana/pyroscope-go"
	"google.golang.org/grpc"
)

func PyroscopeUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var resp interface{}
		var err error

		pyroscope.TagWrapper(ctx, pyroscope.Labels(
			"grpc_method", info.FullMethod,
		), func(ctx context.Context) {
			resp, err = handler(ctx, req)
		})

		return resp, err
	}
}
