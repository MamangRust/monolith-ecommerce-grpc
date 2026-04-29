package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	userhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/user"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	user_cache "github.com/MamangRust/monolith-ecommerce-grpc-user/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	app_errors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	gapi "github.com/MamangRust/monolith-ecommerce-grpc-user/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/service"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
)

type UserHandlerTestSuite struct {
	tests.BaseTestSuite
	client      pb.UserCommandServiceClient
	router      *echo.Echo
	userID      int
	userEmail   string
}

func (s *UserHandlerTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	s.SetupRoleService()
	roleClient := pb.NewRoleQueryServiceClient(s.Conns["role"])

	queries := db.New(s.DBPool())
	repos := repository.NewRepositories(queries, roleClient)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheStore := s.GetCacheStore()
	mencache := user_cache.NewMencache(cacheStore)

	userService := service.NewService(&service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Hash:          hasher,
		Logger:        log,
		Observability: s.Obs,
	})

	// Start gRPC Server
	userHandler := gapi.NewHandler(&gapi.Deps{
		Service: userService,
		Logger:  log,
	})
	server := grpc.NewServer()
	pb.RegisterUserQueryServiceServer(server, userHandler.UserQuery)
	pb.RegisterUserCommandServiceServer(server, userHandler.UserCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)
	s.client = pb.NewUserCommandServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	userhandler.RegisterUserHandler(&userhandler.DepsUser{
		Client:     conn,
		E:          s.router,
		Logger:     log,
		Cache:      cacheStore,
		ApiHandler: app_errors.NewApiHandler(s.Obs, log),
	})
}

func (s *UserHandlerTestSuite) TestUserApiLifecycle() {
	// 1. Create
	s.userEmail = "handler.user@example.com"
	createReq := requests.CreateUserRequest{
		FirstName:       "Handler",
		LastName:        "User",
		Email:           s.userEmail,
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/user-command/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.userID = int(data["id"].(float64))

	// 2. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/user-query", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. FindById
	s.Require().NotZero(s.userID)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/user-query/%d", s.userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/user-query/active", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Update
	updateReq := requests.UpdateUserRequest{
		FirstName:       "Updated",
		LastName:        "User",
		Email:           s.userEmail,
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/user-command/update/%d", s.userID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/user-command/trashed/%d", s.userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/user-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/user-command/restore/%d", s.userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/user-command/trashed/%d", s.userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/user-command/permanent/%d", s.userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/user-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/user-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestUserHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserHandlerTestSuite))
}
