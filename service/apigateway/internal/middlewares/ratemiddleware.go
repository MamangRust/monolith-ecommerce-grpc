package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter
}

func NewRateLimiter(rps int, burst int) *RateLimiter {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	return &RateLimiter{
		limiter: limiter,
	}
}

func (rl *RateLimiter) Limit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !rl.limiter.Allow() {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Too many requests, please try again later",
			})
		}
		return next(c)
	}
}
