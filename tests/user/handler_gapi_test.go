package user_test

import (
	"context"
	"testing"

	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	pbuserrole "github.com/MamangRust/monolith-ecommerce-pb/user_role"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	gapi "github.com/MamangRust/monolith-ecommerce-user/handler"
	"github.com/MamangRust/monolith-ecommerce-user/repository"
	"github.com/MamangRust/monolith-ecommerce-user/service"

	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	user_cache "github.com/MamangRust/monolith-ecommerce-user/cache"
)

type UserGapiTestSuite struct {
	tests.BaseTestSuite
	client      pbuser.UserCommandServiceClient
	queryClient pbuser.UserQueryServiceClient
	userID      int
}

func (s *UserGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	s.SetupRoleService()
	roleClient := pbrole.NewRoleQueryServiceClient(s.Conns["role"])

	queries := db.New(s.DBPool())
	repos := repository.NewRepositories(&repository.Deps{
		Db:              queries,
		RoleQueryClient: roleClient,
		UserRoleClient:  pbuserrole.NewUserRoleServiceClient(s.Conns["role"]),
	})

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheStore := s.GetCacheStore()
	mencache := user_cache.NewMencache(cacheStore)

	userService := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Hash:          hasher,
		Cache:         mencache,
		Observability: s.Obs,
	})

	// Start gRPC Server
	userHandler := gapi.NewHandler(&gapi.Deps{
		Service: userService,
		Logger:  log,
	})
	server := grpc.NewServer()
	pbuser.RegisterUserCommandServiceServer(server, userHandler.UserCommand)
	pbuser.RegisterUserQueryServiceServer(server, userHandler.UserQuery)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.client = pbuser.NewUserCommandServiceClient(conn)
	s.queryClient = pbuser.NewUserQueryServiceClient(conn)
}

func (s *UserGapiTestSuite) TestUserGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pbuser.CreateUserRequest{
		Firstname:       "Gapi",
		Lastname:        "User",
		Email:           "gapi.user@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	res, err := s.client.Create(ctx, createReq)
	s.Require().NoError(err)
	s.Equal(createReq.Email, res.Data.Email)
	userID := res.Data.Id

	// 2. FindById
	getRes, err := s.queryClient.FindById(ctx, &pbuser.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)
	s.Equal(userID, getRes.Data.Id)

	// 3. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbuser.FindAllUserRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbuser.FindAllUserRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, err := s.client.Update(ctx, &pbuser.UpdateUserRequest{
		Id:              userID,
		Firstname:       "GapiUpdated",
		Lastname:        "UserUpdated",
		Email:           "gapi.updated@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	s.Require().NoError(err)
	s.Equal("GapiUpdated", updateRes.Data.Firstname)

	// 6. Trash
	_, err = s.client.TrashedUser(ctx, &pbuser.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbuser.FindAllUserRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.client.RestoreUser(ctx, &pbuser.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, _ = s.client.TrashedUser(ctx, &pbuser.FindByIdUserRequest{Id: userID})
	_, err = s.client.DeleteUserPermanent(ctx, &pbuser.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)

	// 10. RestoreAll
	_, err = s.client.RestoreAllUser(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 11. DeleteAll
	_, err = s.client.DeleteAllUserPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestUserGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserGapiTestSuite))
}
