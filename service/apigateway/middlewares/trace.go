package middlewares

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

// TraceMiddleware extracts trace context from incoming HTTP headers and starts
// a server span for each request, enabling end-to-end distributed tracing from
// the API gateway down to the gRPC dependency services.
func TraceMiddleware(serviceName string) echo.MiddlewareFunc {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			ctx := propagator.Extract(req.Context(), propagation.HeaderCarrier(req.Header))

			route := c.Path()
			if route == "" {
				route = "unmatched"
			}
			spanName := req.Method + " " + route

			ctx, span := tracer.Start(ctx, spanName)
			defer span.End()

			span.SetAttributes(
				attribute.String("http.method", req.Method),
				attribute.String("http.route", route),
				attribute.String("http.target", req.URL.Path),
				attribute.String("http.host", req.Host),
			)

			// Propagate the traced context into the request so downstream
			// gRPC client interceptors can continue the same trace.
			c.SetRequest(req.WithContext(ctx))

			err := next(c)

			span.SetAttributes(attribute.Int("http.status_code", c.Response().Status))
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
			} else if c.Response().Status >= 500 {
				span.SetStatus(codes.Error, "server error")
			}

			return err
		}
	}
}
