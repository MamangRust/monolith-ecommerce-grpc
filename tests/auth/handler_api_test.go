package auth_test

import (
	"bytes"
	"context"
	"fmt"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	authhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/auth"
	auth_cache_api "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/auth"
	"github.com/MamangRust/monolith-ecommerce-auth/handler"
	"github.com/MamangRust/monolith-ecommerce-auth/repository"
	"github.com/MamangRust/monolith-ecommerce-auth/service"
	auth_cache "github.com/MamangRust/monolith-ecommerce-auth/cache"
	user_handler "github.com/MamangRust/monolith-ecommerce-grpc-user/handler"
	user_service "github.com/MamangRust/monolith-ecommerce-grpc-user/service"
	user_repo "github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	role_handler "github.com/MamangRust/monolith-ecommerce-grpc-role/handler"
	role_service "github.com/MamangRust/monolith-ecommerce-grpc-role/service"
	role_repo "github.com/MamangRust/monolith-ecommerce-grpc-role/repository"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	user_cache "github.com/MamangRust/monolith-ecommerce-grpc-user/cache"
	role_cache "github.com/MamangRust/monolith-ecommerce-grpc-role/cache"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type AuthHandlerApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	server      *echo.Echo
	email       string
	password    string
	accessToken string
	userID      int
}

func (s *AuthHandlerApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)
	s.redisClient.FlushAll(context.Background())

	queries := db.New(pool)
	
	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	obs, _ := observability.NewObservability("test", log)

	// 1. Setup Role Service & gRPC Server
	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(queries)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repository:    roleRepos,
		Logger:        log,
		Cache:         roleMencache,
		Observability: obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  log,
	})
	roleServer := grpc.NewServer()
	pb.RegisterRoleQueryServiceServer(roleServer, roleGapi.RoleQuery)
	pb.RegisterRoleCommandServiceServer(roleServer, roleGapi.RoleCommand)
	roleLis, _ := net.Listen("tcp", "localhost:0")
	go roleServer.Serve(roleLis)
	roleConn, _ := grpc.NewClient(roleLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 2. Setup User Service & gRPC Server
	userMencache := user_cache.NewMencache(cacheStore)
	roleQueryClientForUser := pb.NewRoleQueryServiceClient(roleConn)
	userRepos := user_repo.NewRepositories(queries, roleQueryClientForUser)
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        log,
		Hash:         hasher,
		Cache:         userMencache,
		Observability: obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  log,
	})
	userServer := grpc.NewServer()
	pb.RegisterUserQueryServiceServer(userServer, userGapi.UserQuery)
	pb.RegisterUserCommandServiceServer(userServer, userGapi.UserCommand)
	userLis, _ := net.Listen("tcp", "localhost:0")
	go userServer.Serve(userLis)
	userConn, _ := grpc.NewClient(userLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 3. Setup Auth Service with gRPC clients
	userQueryClient := pb.NewUserQueryServiceClient(userConn)
	userCommandClient := pb.NewUserCommandServiceClient(userConn)
	roleQueryClient := pb.NewRoleQueryServiceClient(roleConn)
	roleCommandClient := pb.NewRoleCommandServiceClient(roleConn)

	repos := repository.NewRepositories(queries, userQueryClient, userCommandClient, roleQueryClient, roleCommandClient)
	
	tokenManager, _ := auth.NewManager("mysecret")
	mencache := auth_cache.NewMencache(cacheStore)
	apiAuthCache := auth_cache_api.NewMencache(cacheStore)

	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: obs,
	})

	h := handler.NewAuthHandleGrpc(svc, log)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, h)

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.server = echo.New()
	apiHandler := errors.NewApiHandler(obs, log)

	// Auth bypass middleware for /api/auth/me
	s.server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if s.userID != 0 {
				c.Set("userId", strconv.Itoa(s.userID))
			}
			return next(c)
		}
	})

	fmt.Println("Calling RegisterAuthHandler...")
	authhandler.RegisterAuthHandler(&authhandler.DepsAuth{
		Client:     conn,
		E:          s.server,
		Logger:     log,
		Cache:      apiAuthCache,
		ApiHandler: apiHandler,
	})
	fmt.Println("RegisterAuthHandler called.")

	s.server.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	s.email = "auth.handler.api.test@example.com"
	s.password = "password123"

	// Seed ROLE_ADMIN via gRPC to ensure visibility
	roleCommandClient = pb.NewRoleCommandServiceClient(roleConn)
	createdRoleRes, err := roleCommandClient.CreateRole(context.Background(), &pb.CreateRoleRequest{
		Name: "ROLE_ADMIN",
	})
	if err != nil {
		fmt.Printf("DEBUG: Seed ROLE_ADMIN gRPC error: %v\n", err)
	} else {
		fmt.Printf("DEBUG: Seed ROLE_ADMIN gRPC success, ID: %d\n", createdRoleRes.Data.Id)
	}

	// Verify via direct queries
	testRes, err := queries.GetRoles(context.Background(), db.GetRolesParams{
		Column1: "",
		Limit:   10,
		Offset:  0,
	})
	if err != nil {
		fmt.Printf("DEBUG: direct queries error: %v\n", err)
	} else {
		fmt.Printf("DEBUG: direct queries found %d roles\n", len(testRes))
	}

	// Verify via GetRole by ID
	roleByID, err := queries.GetRole(context.Background(), createdRoleRes.Data.Id)
	if err != nil {
		fmt.Printf("DEBUG: GetRole by ID error: %v\n", err)
	} else {
		fmt.Printf("DEBUG: GetRole by ID found: %s\n", roleByID.RoleName)
	}

	// Verify via gRPC with empty search
	roleClient := pb.NewRoleQueryServiceClient(roleConn)
	roleRes, err := roleClient.FindAllRole(context.Background(), &pb.FindAllRoleRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		fmt.Printf("DEBUG: gRPC verify error: %v\n", err)
	} else {
		fmt.Printf("DEBUG: gRPC verify (empty search) found %d roles\n", len(roleRes.Data))
		for _, r := range roleRes.Data {
			fmt.Printf("- gRPC role: %s\n", r.Name)
		}
	}
}

func (s *AuthHandlerApiTestSuite) TearDownSuite() {
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	s.ts.Teardown()
}

func (s *AuthHandlerApiTestSuite) Test0_Ping() {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal("pong", rec.Body.String())
}

func (s *AuthHandlerApiTestSuite) Test0_Hello() {
	req := httptest.NewRequest(http.MethodGet, "/api/auth/hello", nil)
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal("Hello", rec.Body.String())
}

func (s *AuthHandlerApiTestSuite) Test1_Register() {
	body := map[string]string{
		"firstname":        "Auth",
		"lastname":         "API",
		"email":            s.email,
		"password":         s.password,
		"confirm_password": s.password,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.server.ServeHTTP(rec, req)

	s.Equal(http.StatusCreated, rec.Code, "Expected StatusCreated, got %d. Body: %s", rec.Code, rec.Body.String())
	
	var res map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)
	
	data, ok := res["data"].(map[string]interface{})
	s.True(ok, "Expected 'data' to be a map, got %T. Body: %s", res["data"], rec.Body.String())
	s.userID = int(data["id"].(float64))
}

func (s *AuthHandlerApiTestSuite) Test2_Login() {
	body := map[string]string{
		"email":    s.email,
		"password": s.password,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.server.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code, "Expected StatusOK, got %d. Body: %s", rec.Code, rec.Body.String())
	
	var res map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)
	
	data, ok := res["data"].(map[string]interface{})
	s.True(ok, "Expected 'data' to be a map, got %T. Body: %s", res["data"], rec.Body.String())
	s.accessToken = data["access_token"].(string)
}

func (s *AuthHandlerApiTestSuite) Test4_LoginLockout() {
	email := "locked.api@example.com"
	password := "wrongpassword"

	// Register user first
	regBody := map[string]string{
		"firstname":        "Locked",
		"lastname":         "API",
		"email":            email,
		"password":         "correctpassword",
		"confirm_password": "correctpassword",
	}
	jsonRegBody, _ := json.Marshal(regBody)
	regReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonRegBody))
	regReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	regRec := httptest.NewRecorder()
	s.server.ServeHTTP(regRec, regReq)
	s.Equal(http.StatusCreated, regRec.Code)

	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonLoginBody, _ := json.Marshal(loginBody)

	// Fail login 5 times (total 5)
	for i := 0; i < 5; i++ {
		loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonLoginBody))
		loginReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		loginRec := httptest.NewRecorder()
		s.server.ServeHTTP(loginRec, loginReq)
		s.Equal(http.StatusUnauthorized, loginRec.Code)
	}

	// 6th attempt should return 403 Forbidden (ErrAccountLocked)
	lockedReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonLoginBody))
	lockedReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	lockedRec := httptest.NewRecorder()
	s.server.ServeHTTP(lockedRec, lockedReq)
	s.Equal(http.StatusForbidden, lockedRec.Code)
}

func (s *AuthHandlerApiTestSuite) Test3_GetMe() {
	s.Require().NotZero(s.userID)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()

	s.server.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code, "Expected StatusOK, got %d. Body: %s", rec.Code, rec.Body.String())
	
	var res map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)

	data, ok := res["data"].(map[string]interface{})
	s.True(ok, "Expected 'data' to be a map, got %T. Body: %s", res["data"], rec.Body.String())
	s.Equal(s.email, data["email"])
}

func TestAuthHandlerApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthHandlerApiTestSuite))
}
