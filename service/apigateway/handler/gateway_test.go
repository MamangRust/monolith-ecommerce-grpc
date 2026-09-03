package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/middlewares"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	echoSwagger "github.com/swaggo/echo-swagger"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func testLogger(t *testing.T) logger.LoggerInterface {
	logger.ResetInstance()
	l, err := logger.NewLogger("apigateway-test", sdklog.NewLoggerProvider())
	require.NoError(t, err)
	require.NotNil(t, l)
	return l
}

// bootGateway mirrors the production wiring (apps/client.go createEchoServer +
// handler.NewHandler): global JWT middleware with the real whitelist, swagger,
// health, and every registered domain handler. gRPC connections may be nil
// because handlers only use them at request time.
func bootGateway(t *testing.T, conns *handler.ServiceConnections) *echo.Echo {
	viper.Set("SECRET_KEY", "test-secret")
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	middlewares.RegisterErrorHandler(e)
	middlewares.WebSecurityConfig(e)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})
	handler.NewHandler(&handler.Deps{
		E:                  e,
		Logger:             testLogger(t),
		ServiceConnections: conns,
		Cache:              nil,
		Image:              nil,
		Kafka:              nil,
	})
	return e
}

func hasRoute(routes []*echo.Route, method, path string) bool {
	for _, r := range routes {
		if r.Method == method && r.Path == path {
			return true
		}
	}
	return false
}

func signToken(t *testing.T, sub string) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": sub})
	signed, err := tok.SignedString([]byte("test-secret"))
	require.NoError(t, err)
	return signed
}

func perform(e *echo.Echo, method, path string, body io.Reader, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	if token != "" {
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// TestGatewayRouteInventory validates that the live route table, built from the
// actual handler registration, follows the -command/-query contract and has no
// swapped groups (regression for the cart route swap).
func TestGatewayRouteInventory(t *testing.T) {
	e := bootGateway(t, &handler.ServiceConnections{})

	routes := e.Routes()
	require.NotEmpty(t, routes, "gateway must register routes")

	// Cart command/query separation (regression for the swapped groups).
	require.True(t, hasRoute(routes, http.MethodPost, "/api/cart-command/create"))
	require.True(t, hasRoute(routes, http.MethodDelete, "/api/cart-command/delete"))
	require.True(t, hasRoute(routes, http.MethodPost, "/api/cart-command/delete-all"))
	require.True(t, hasRoute(routes, http.MethodGet, "/api/cart-query"))
	require.False(t, hasRoute(routes, http.MethodPost, "/api/cart-query/create"),
		"command must not be registered under the cart-query group")

	// Convention: groups ending in -command only carry mutations; -query only GETs.
	for _, r := range routes {
		api := strings.TrimPrefix(r.Path, "/api/")
		group := strings.SplitN(api, "/", 2)[0]
		switch {
		case strings.HasSuffix(group, "-command"):
			require.NotEqual(t, http.MethodGet, r.Method, "GET must not be a command route: %s", r.Path)
		case strings.HasSuffix(group, "-query"):
			require.Equal(t, http.MethodGet, r.Method, "only GET allowed on query group: %s", r.Path)
		}
	}

	require.True(t, hasRoute(routes, http.MethodGet, "/health"))
	require.True(t, hasRoute(routes, http.MethodGet, "/swagger/*"))
}

// TestSwaggerMatchesRegisteredRoutes enforces that the generated Swagger spec
// documents every route actually registered by the gateway (and nothing else of
// significance). Echo :params are normalized to Swagger {params} before the set
// comparison; regenerate with `just generate-swagger` after route changes.
func TestSwaggerMatchesRegisteredRoutes(t *testing.T) {
	e := bootGateway(t, &handler.ServiceConnections{})

	raw, err := os.ReadFile("../docs/swagger.json")
	require.NoError(t, err, "run `just generate-swagger` to produce docs/swagger.json")

	var spec struct {
		Paths map[string]json.RawMessage `json:"paths"`
	}
	require.NoError(t, json.Unmarshal(raw, &spec))

	paramRegex := regexp.MustCompile(`:([a-zA-Z_]+)`)
	var missing []string
	for _, r := range e.Routes() {
		if !strings.HasPrefix(r.Path, "/api/") {
			continue
		}
		swaggerPath := paramRegex.ReplaceAllString(r.Path, "{$1}")
		if _, ok := spec.Paths[swaggerPath]; !ok {
			missing = append(missing, r.Method+" "+r.Path)
		}
	}
	require.Empty(t, missing,
		"registered routes missing from swagger.json (regenerate docs and fix annotations):\n%s",
		strings.Join(missing, "\n"))
}

// TestGatewayMiddlewareSmoke covers the exit-criteria status matrix on the real
// middleware chain: health 200 without token, 401, 404, 400 (validation), and
// 503 (unavailable dependency).
func TestGatewayMiddlewareSmoke(t *testing.T) {
	deadConn, err := grpc.NewClient("localhost:1", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	conns := &handler.ServiceConnections{Cart: deadConn}
	e := bootGateway(t, conns)
	token := signToken(t, "42")

	// Health must be reachable without a token (JWT whitelist).
	rec := perform(e, http.MethodGet, "/health", nil, "")
	require.Equal(t, http.StatusOK, rec.Code, "health must not require auth")

	// Protected route without a token → 401.
	rec = perform(e, http.MethodGet, "/api/cart-query", nil, "")
	require.Equal(t, http.StatusUnauthorized, rec.Code)

	// Unknown route with a valid token → 404.
	rec = perform(e, http.MethodGet, "/api/does-not-exist", nil, token)
	require.Equal(t, http.StatusNotFound, rec.Code)

	// Invalid JSON body with a valid token → 400 (before the gRPC call).
	rec = perform(e, http.MethodPost, "/api/cart-command/create", strings.NewReader("{not-json"), token)
	require.Equal(t, http.StatusBadRequest, rec.Code)

	// Valid token + unavailable downstream → 503.
	rec = perform(e, http.MethodPost, "/api/cart-command/create", strings.NewReader(`{"product_id":1,"quantity":1}`), token)
	require.Equal(t, http.StatusServiceUnavailable, rec.Code, "unavailable gRPC dependency must map to 503")
}
