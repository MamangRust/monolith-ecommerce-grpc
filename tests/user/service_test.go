package user_test

import (
	"context"
	"testing"

	user_cache "github.com/MamangRust/monolith-ecommerce-grpc-user/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"

	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type UserServiceTestSuite struct {
	tests.BaseTestSuite
	userID      int
	userService *service.Service
}

func (s *UserServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	queries := db.New(s.DBPool())

	// Setup Role connection for repository
	s.SetupRoleService()
	roleClient := pb.NewRoleQueryServiceClient(s.Conns["role"])

	repos := repository.NewRepositories(queries, roleClient)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheStore := s.GetCacheStore()
	mencache := user_cache.NewMencache(cacheStore)

	s.userService = service.NewService(&service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Hash:          hasher,
		Logger:        log,
		Observability: s.Obs,
	})
}

func (s *UserServiceTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *UserServiceTestSuite) TestUserLifecycle() {
	ctx := context.Background()

	// 1. Setup Dependencies
	s.SetupRoleService()

	// 2. Create User
	req := &requests.CreateUserRequest{
		FirstName:       "Test",
		LastName:        "User",
		Email:           "testuser_lifecycle@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	created, err := s.userService.UserCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	userID := int(created.UserID)

	// 3. FindByID
	found, err := s.userService.UserQuery.FindByID(ctx, userID)
	s.Require().NoError(err)
	s.Equal(req.Email, found.Email)

	// 4. FindByEmailWithPassword
	foundByEmail, err := s.userService.UserQuery.FindByEmailWithPassword(ctx, req.Email)
	s.Require().NoError(err)
	s.Equal(userID, int(foundByEmail.UserID))

	// 5. Update
	newFirstname := "UpdatedTest"
	updateReq := &requests.UpdateUserRequest{
		UserID:          &userID,
		FirstName:       newFirstname,
		LastName:        req.LastName,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}
	updated, err := s.userService.UserCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(newFirstname, updated.Firstname)

	// 6. FindAll
	_, total, err := s.userService.UserQuery.FindAll(ctx, &requests.FindAllUsers{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 7. Trash
	_, err = s.userService.UserCommand.Trash(ctx, userID)
	s.Require().NoError(err)

	// 8. FindTrashed
	_, totalTrashed, err := s.userService.UserQuery.FindTrashed(ctx, &requests.FindAllUsers{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 9. FindActive
	active, _, err := s.userService.UserQuery.FindActive(ctx, &requests.FindAllUsers{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, u := range active {
		s.NotEqual(userID, int(u.UserID))
	}

	// 10. Restore
	_, err = s.userService.UserCommand.Restore(ctx, userID)
	s.Require().NoError(err)

	// 11. DeletePermanent
	_, err = s.userService.UserCommand.Trash(ctx, userID)
	s.Require().NoError(err)
	success, err := s.userService.UserCommand.DeletePermanent(ctx, userID)
	s.Require().NoError(err)
	s.True(success)

	// 12. RestoreAll & DeleteAll
	u1, _ := s.userService.UserCommand.Create(ctx, &requests.CreateUserRequest{FirstName: "U1", LastName: "L1", Email: "u1@x.com", Password: "p1", ConfirmPassword: "p1"})
	u2, _ := s.userService.UserCommand.Create(ctx, &requests.CreateUserRequest{FirstName: "U2", LastName: "L2", Email: "u2@x.com", Password: "p2", ConfirmPassword: "p2"})
	
	s.userService.UserCommand.Trash(ctx, int(u1.UserID))
	s.userService.UserCommand.Trash(ctx, int(u2.UserID))

	resRestoreAll, err := s.userService.UserCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.userService.UserCommand.Trash(ctx, int(u1.UserID))
	s.userService.UserCommand.Trash(ctx, int(u2.UserID))

	resDeleteAll, err := s.userService.UserCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func (s *UserServiceTestSuite) Test1_CreateUser() {
	ctx := context.Background()

	req := &requests.CreateUserRequest{
		FirstName: "User",
		LastName:  "Service",
		Email:     "user.service@example.com",
		Password:  "password123",
	}
	user, err := s.userService.UserCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(user)
	s.Equal(req.Email, user.Email)
	s.userID = int(user.UserID)
}

func (s *UserServiceTestSuite) Test2_FindUserById() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	found, err := s.userService.UserQuery.FindByID(ctx, s.userID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.userID, int(found.UserID))
}

func (s *UserServiceTestSuite) Test3_UpdateUser() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	updateReq := &requests.UpdateUserRequest{
		UserID:    &s.userID,
		FirstName: "Updated",
	}
	updated, err := s.userService.UserCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated", updated.Firstname)
}

func (s *UserServiceTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	_, err := s.userService.UserCommand.Trash(ctx, s.userID)
	s.NoError(err)

	_, err = s.userService.UserCommand.Restore(ctx, s.userID)
	s.NoError(err)
}

func (s *UserServiceTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	success, err := s.userService.UserCommand.DeletePermanent(ctx, s.userID)
	s.NoError(err)
	s.True(success)
}

func TestUserServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserServiceTestSuite))
}
