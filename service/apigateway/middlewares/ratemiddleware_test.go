package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestRateLimiterAllowsWithinBurst(t *testing.T) {
	rl := NewRateLimiter(1, 2) // 1 req/sec with burst of 2
	e := echo.New()
	handler := rl.Limit(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		require.NoError(t, handler(e.NewContext(req, rec)))
		require.Equal(t, http.StatusOK, rec.Code, "request %d must be allowed within burst", i+1)
	}
}

func TestRateLimiterReturns429BeyondBurst(t *testing.T) {
	rl := NewRateLimiter(1, 2) // 1 req/sec with burst of 2
	e := echo.New()
	handler := rl.Limit(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		require.NoError(t, handler(e.NewContext(req, rec)))
	}

	// Third request exceeds the burst → 429.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	require.NoError(t, handler(e.NewContext(req, rec)))
	require.Equal(t, http.StatusTooManyRequests, rec.Code)
	require.Contains(t, rec.Body.String(), "Too many requests")
}
