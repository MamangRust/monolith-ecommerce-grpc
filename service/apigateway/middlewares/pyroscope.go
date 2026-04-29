package middlewares

import (
	"context"
	"github.com/grafana/pyroscope-go"
	"github.com/labstack/echo/v4"
)

func PyroscopeMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			pyroscope.TagWrapper(c.Request().Context(), pyroscope.Labels(
				"path", c.Path(),
				"method", c.Request().Method,
			), func(ctx context.Context) {
				err = next(c)
			})
			return err
		}
	}
}
