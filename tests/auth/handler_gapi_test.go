package auth_test

import (
	"context"
	"net"
	"strconv"
	"testing"

	mencache "github.com/MamangRust/monolith-ecommerce-auth/cache"
	"github.com/MamangRust/monolith-ecommerce-auth/handler"
	"github.com/MamangRust/monolith-ecommerce-auth/repository"
	"github.com/MamangRust/monolith-ecommerce-auth/service"
	pb "github.com/MamangRust/monolith-ecommerce-pb"
	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	pbuserrole "github.com/MamangRust/monolith-ecommerce-pb/user_role"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"

	role_cache "github.com/MamangRust/monolith-ecommerce-role/cache"
	role_handler "github.com/MamangRust/monolith-ecommerce-role/handler"
	role_repo "github.com/MamangRust/monolith-ecommerce-role/repository"
	role_service "github.com/MamangRust/monolith-ecommerce-role/service"
	user_cache "github.com/MamangRust/monolith-ecommerce-user/cache"
	user_handler "github.com/MamangRust/monolith-ecommerce-user/handler"
	user_repo "github.com/MamangRust/monolith-ecommerce-user/repository"
	user_service "github.com/MamangRust/monolith-ecommerce-user/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthHandlerGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	client      pb.AuthServiceClient
	conn        *grpc.ClientConn
	grpcServer  *grpc.Server
	email       string
	password    string
	accessToken string
}

func (s *AuthHandlerGapiTestSuite) SetupSuite() {
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
	pbrole.RegisterRoleQueryServiceServer(roleServer, roleGapi.RoleQuery)
	pbrole.RegisterRoleCommandServiceServer(roleServer, roleGapi.RoleCommand)
	pbuserrole.RegisterUserRoleServiceServer(roleServer, roleGapi.UserRole)
	roleLis, _ := net.Listen("tcp", "localhost:0")
	go roleServer.Serve(roleLis)
	roleConn, _ := grpc.NewClient(roleLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 2. Setup User Service & gRPC Server
	userMencache := user_cache.NewMencache(cacheStore)
	roleQueryClientForUser := pbrole.NewRoleQueryServiceClient(roleConn)
	userRepos := user_repo.NewRepositories(&user_repo.Deps{
		Db:              queries,
		RoleQueryClient: roleQueryClientForUser,
		UserRoleClient:  pbuserrole.NewUserRoleServiceClient(roleConn),
	})
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        log,
		Hash:          hasher,
		Cache:         userMencache,
		Observability: obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  log,
	})
	userServer := grpc.NewServer()
	pbuser.RegisterUserQueryServiceServer(userServer, userGapi.UserQuery)
	pbuser.RegisterUserCommandServiceServer(userServer, userGapi.UserCommand)
	userLis, _ := net.Listen("tcp", "localhost:0")
	go userServer.Serve(userLis)
	userConn, _ := grpc.NewClient(userLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 3. Setup Auth Service with gRPC clients
	userQueryClient := pbuser.NewUserQueryServiceClient(userConn)
	userCommandClient := pbuser.NewUserCommandServiceClient(userConn)
	roleQueryClient := pbrole.NewRoleQueryServiceClient(roleConn)
	roleCommandClient := pbrole.NewRoleCommandServiceClient(roleConn)

	repos := repository.NewRepositories(&repository.Deps{
		Db:                queries,
		UserQueryClient:   userQueryClient,
		UserCommandClient: userCommandClient,
		RoleQueryClient:   roleQueryClient,
		RoleCommandClient: roleCommandClient,
		UserRoleClient:    pbuserrole.NewUserRoleServiceClient(roleConn),
	})

	tokenManager, _ := auth.NewManager("mysecret")
	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache.NewMencache(cacheStore),
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: obs,
	})

	h := handler.NewAuthHandleGrpc(svc, log)

	s.grpcServer = grpc.NewServer()
	pb.RegisterAuthServiceServer(s.grpcServer, h)

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	go func() {
		_ = s.grpcServer.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewAuthServiceClient(conn)

	s.email = "auth.handler.gapi.test@example.com"
	s.password = "password123"

	// Seed ROLE_ADMIN
	_, _ = pool.Exec(context.Background(), "INSERT INTO roles (role_name) VALUES ('ROLE_ADMIN')")
}

func (s *AuthHandlerGapiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	s.ts.Teardown()
}

func (s *AuthHandlerGapiTestSuite) Test1_Register() {
	ctx := context.Background()
	req := &pb.RegisterRequest{
		Firstname:       "Auth",
		Lastname:        "Handler",
		Email:           s.email,
		Password:        s.password,
		ConfirmPassword: s.password,
	}

	res, err := s.client.RegisterUser(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("success", res.Status)
	s.Equal(s.email, res.Data.Email)
}

func (s *AuthHandlerGapiTestSuite) Test2_Login() {
	ctx := context.Background()
	req := &pb.LoginRequest{
		Email:    s.email,
		Password: s.password,
	}

	res, err := s.client.LoginUser(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data.AccessToken)
	s.accessToken = res.Data.AccessToken
}

func (s *AuthHandlerGapiTestSuite) Test4_LoginLockout() {
	ctx := context.Background()
	email := "locked.gapi@example.com"
	password := "wrongpassword"

	// Register user first
	regReq := &pb.RegisterRequest{
		Firstname:       "Locked",
		Lastname:        "Gapi",
		Email:           email,
		Password:        "correctpassword",
		ConfirmPassword: "correctpassword",
	}
	_, err := s.client.RegisterUser(ctx, regReq)
	s.NoError(err)

	loginReq := &pb.LoginRequest{
		Email:    email,
		Password: password,
	}

	// Fail login 5 times
	for i := 0; i < 5; i++ {
		_, err := s.client.LoginUser(ctx, loginReq)
		s.Error(err)
	}

	// 6th attempt should return error
	_, err = s.client.LoginUser(ctx, loginReq)
	s.Error(err)
	s.Contains(err.Error(), "Account temporarily locked")
}

func (s *AuthHandlerGapiTestSuite) Test3_GetMe() {
	s.Require().NotEmpty(s.accessToken)
	ctx := context.Background()

	tokenManager, _ := auth.NewManager("mysecret")
	userIdStr, err := tokenManager.ValidateToken(s.accessToken)
	s.NoError(err)

	userId, err := strconv.Atoi(userIdStr)
	s.NoError(err)

	res, err := s.client.GetMe(ctx, &pb.GetMeRequest{UserId: int32(userId)})
	s.NoError(err)
	s.NotNil(res)
	s.Equal("success", res.Status)
	s.Equal(s.email, res.Data.Email)
}

func TestAuthHandlerGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthHandlerGapiTestSuite))
}
