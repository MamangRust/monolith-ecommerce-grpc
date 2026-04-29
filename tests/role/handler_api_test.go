package role_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	rolehandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/role"
	apicache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache"
	role_cache "github.com/MamangRust/monolith-ecommerce-grpc-role/cache"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	app_errors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RoleApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	echo        *echo.Echo
	grpcServer  *grpc.Server
	conn        *grpc.ClientConn
	roleID      int
}

func (s *RoleApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	queries := db.New(pool)
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	mencache := role_cache.NewMencache(cacheStore)

	roleService := service.NewService(&service.Deps{
		Repository:    repos,
		Logger:        log,
		Cache:         mencache,
		Observability: obs,
	})

	// Start internal gRPC Server for Role module
	roleHandlerGrpc := handler.NewHandler(&handler.Deps{
		Service: roleService,
		Logger:  log,
	})
	server := grpc.NewServer()
	pb.RegisterRoleCommandServiceServer(server, roleHandlerGrpc.RoleCommand)
	pb.RegisterRoleQueryServiceServer(server, roleHandlerGrpc.RoleQuery)
	s.grpcServer = server

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn

	// Setup Echo and API Handler
	e := echo.New()
	s.echo = e

	// Bypass auth middleware by setting user_id and seeding roles in Redis
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user_id", 1)
			return next(c)
		}
	})

	roles := []string{"Admin_Role_10", "Admin_Admin_14"}
	cache.SetToCache(s.ts.Ctx, cacheStore, "user_roles:1", &roles, 5*time.Minute)

	apiErrorHandler := app_errors.NewApiHandler(obs, log)
	rolehandler.RegisterRoleHandler(&rolehandler.DepsRole{
		Client:     conn,
		Kafka:      nil,
		E:          e,
		Logger:     log,
		CacheStore: cacheStore,
		Cache:      apicache.NewRoleCache(cacheStore),
		ApiHandler: apiErrorHandler,
	})
}

func (s *RoleApiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *RoleApiTestSuite) TestRoleApiLifecycle() {
	// 1. Create
	reqBody := requests.CreateRoleRequest{
		Name: "API Role",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/role-command/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal(reqBody.Name, data["name"])
	s.roleID = int(data["id"].(float64))

	// 2. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/role-query", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. FindById
	s.Require().NotZero(s.roleID)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/role-query/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.roleID), data["id"])

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. FindByUserId
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/user/1", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Update
	updateBody := requests.UpdateRoleRequest{
		Name: "Updated API Role",
	}
	body, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role-command/update/%d", s.roleID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role-command/restore/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeletePermanent
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/role-command/permanent/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/role-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/role-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestRoleApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleApiTestSuite))
}
