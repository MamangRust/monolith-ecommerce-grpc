package role_test

import (
	"context"
	"testing"

	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	role_cache "github.com/MamangRust/monolith-ecommerce-role/cache"
	role_handler "github.com/MamangRust/monolith-ecommerce-role/handler"
	"github.com/MamangRust/monolith-ecommerce-role/repository"
	"github.com/MamangRust/monolith-ecommerce-role/service"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RoleGapiTestSuite struct {
	tests.BaseTestSuite
	commandClient pbrole.RoleCommandServiceClient
	queryClient   pbrole.RoleQueryServiceClient
}

func (s *RoleGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Role dependencies
	mencache := role_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(queries)
	svc := service.NewService(&service.Deps{
		Repository:    repos,
		Logger:        s.Log,
		Cache:         mencache,
		Observability: s.Obs,
	})

	// Handler
	handler := role_handler.NewHandler(&role_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbrole.RegisterRoleCommandServiceServer(server, handler.RoleCommand)
	pbrole.RegisterRoleQueryServiceServer(server, handler.RoleQuery)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.commandClient = pbrole.NewRoleCommandServiceClient(conn)
	s.queryClient = pbrole.NewRoleQueryServiceClient(conn)
}

func (s *RoleGapiTestSuite) TestRoleGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createRes, err := s.commandClient.CreateRole(ctx, &pbrole.CreateRoleRequest{
		Name: "Gapi Role",
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	roleID := createRes.Data.Id

	// 2. FindById
	getRes, err := s.queryClient.FindByIdRole(ctx, &pbrole.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)
	s.Equal("Gapi Role", getRes.Data.Name)

	// 3. FindAll
	allRes, err := s.queryClient.FindAllRole(ctx, &pbrole.FindAllRoleRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbrole.FindAllRoleRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, err := s.commandClient.UpdateRole(ctx, &pbrole.UpdateRoleRequest{
		Id:   roleID,
		Name: "Gapi Role Updated",
	})
	s.Require().NoError(err)
	s.Equal("Gapi Role Updated", updateRes.Data.Name)

	// 6. Trash
	_, err = s.commandClient.TrashedRole(ctx, &pbrole.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbrole.FindAllRoleRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.commandClient.RestoreRole(ctx, &pbrole.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, _ = s.commandClient.TrashedRole(ctx, &pbrole.FindByIdRoleRequest{RoleId: roleID})
	_, err = s.commandClient.DeleteRolePermanent(ctx, &pbrole.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)

	// 10. RestoreAll
	_, err = s.commandClient.RestoreAllRole(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 11. DeleteAll
	_, err = s.commandClient.DeleteAllRolePermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestRoleGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleGapiTestSuite))
}
