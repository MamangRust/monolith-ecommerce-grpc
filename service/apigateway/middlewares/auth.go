package middlewares

import (
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

var whiteListPaths = []string{
	"/api/auth/login",
	"/api/auth/register",
	"/api/auth/hello",
	"/api/auth/verify-code",
	"/docs/",
	"/docs",
	"/swagger",
	"/metrics",
}

func extractUserIDFromClaims(claims jwt.MapClaims) (int, bool) {
	subVal, exists := claims["sub"]
	if !exists {
		return 0, false
	}

	switch v := subVal.(type) {
	case string:
		if v == "" {
			return 0, false
		}
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return id, true
	case float64:
		return int(v), true
	case int:
		return v, true
	default:
		return 0, false
	}
}

func WebSecurityConfig(e *echo.Echo) {
	config := echojwt.Config{
		SigningKey: []byte(viper.GetString("SECRET_KEY")),
		Skipper:    skipAuth,
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)
			if claims, ok := user.Claims.(jwt.MapClaims); ok {
				if id, ok := extractUserIDFromClaims(claims); ok && id > 0 {
					c.Set("user_id", id)
				} else {
					c.Logger().Debugf("Invalid or missing user ID from JWT: %#v", claims["sub"])
				}
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.ErrUnauthorized
		},
	}
	e.Use(echojwt.WithConfig(config))
}

func skipAuth(c echo.Context) bool {
	path := c.Request().URL.Path

	for _, p := range whiteListPaths {
		if path == p ||
			strings.HasPrefix(path, "/swagger") || strings.HasPrefix(path, "/api/auth/verify-code") {
			return true
		}
	}

	return false
}
