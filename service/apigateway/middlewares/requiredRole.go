package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw, ok := c.Get("role_names").([]string)
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "Roles not found")
			}

			for _, r := range raw {
				for _, allowed := range allowedRoles {
					if r == allowed {
						return next(c)
					}
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Role not permitted")
		}
	}
}
