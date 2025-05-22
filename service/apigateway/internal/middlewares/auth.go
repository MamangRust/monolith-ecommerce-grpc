package middlewares

import (
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
	"/docs/",
	"/docs",
	"/swagger",
}

func WebSecurityConfig(e *echo.Echo) {
	config := echojwt.Config{
		SigningKey: []byte(viper.GetString("SECRET_KEY")),
		Skipper:    skipAuth,
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)

			if claims, ok := user.Claims.(jwt.MapClaims); ok {
				subject := claims["sub"]
				c.Set("userID", subject)
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.ErrUnauthorized
		},
	}
	e.Use(echojwt.WithConfig(config))

}

func skipAuth(e echo.Context) bool {
	path := e.Path()

	for _, p := range whiteListPaths {
		if path == p || strings.HasPrefix(path, "/swagger") {
			return true
		}
	}

	return false
}
