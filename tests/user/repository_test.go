package user_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo   *repository.Repositories
	userID int
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	queries := db.New(s.DBPool())
	s.SetupRoleService()
	roleClient := pb.NewRoleQueryServiceClient(s.Conns["role"])
	s.repo = repository.NewRepositories(queries, roleClient)
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *UserRepositoryTestSuite) Test01_CreateUser() {
	ctx := context.Background()

	req := &requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	created, err := s.repo.UserCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.FirstName, created.Firstname)
	s.Equal(req.Email, created.Email)
	s.userID = int(created.UserID)
}

func (s *UserRepositoryTestSuite) Test02_FindById() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	found, err := s.repo.UserQuery.FindByID(ctx, s.userID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.userID, int(found.UserID))
}

func (s *UserRepositoryTestSuite) Test03_FindByEmail() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	found, err := s.repo.UserQuery.FindByEmail(ctx, "john.doe@example.com")
	s.NoError(err)
	s.NotNil(found)
	s.Equal("john.doe@example.com", found.Email)
}

func (s *UserRepositoryTestSuite) Test04_FindByEmailWithPassword() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	found, err := s.repo.UserQuery.FindByEmailWithPassword(ctx, "john.doe@example.com")
	s.NoError(err)
	s.NotNil(found)
	s.NotEmpty(found.Password)
}

func (s *UserRepositoryTestSuite) Test05_FindAllUser() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	all, err := s.repo.UserQuery.FindAll(ctx, &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(all)
}

func (s *UserRepositoryTestSuite) Test06_FindActive() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	active, err := s.repo.UserQuery.FindActive(ctx, &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(active)
}

func (s *UserRepositoryTestSuite) Test07_UpdateUser() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	updateReq := &requests.UpdateUserRequest{
		UserID:          &s.userID,
		FirstName:       "Updated",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}
	updated, err := s.repo.UserCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated", updated.Firstname)
}

func (s *UserRepositoryTestSuite) Test08_TrashAndFindTrashed() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	// Trash
	_, err := s.repo.UserCommand.Trash(ctx, s.userID)
	s.NoError(err)

	// FindTrashed
	trashed, err := s.repo.UserQuery.FindTrashed(ctx, &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashed)

	// FindActive — should NOT include trashed user
	active, err := s.repo.UserQuery.FindActive(ctx, &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	for _, u := range active {
		s.NotEqual(s.userID, int(u.UserID))
	}
}

func (s *UserRepositoryTestSuite) Test09_Restore() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	_, err := s.repo.UserCommand.Restore(ctx, s.userID)
	s.NoError(err)

	// Verify it's back in active
	active, err := s.repo.UserQuery.FindActive(ctx, &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	var found bool
	for _, u := range active {
		if int(u.UserID) == s.userID {
			found = true
			break
		}
	}
	s.True(found, "restored user should appear in active list")
}

func (s *UserRepositoryTestSuite) Test10_DeletePermanent() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	// Must trash first
	_, err := s.repo.UserCommand.Trash(ctx, s.userID)
	s.Require().NoError(err)

	success, err := s.repo.UserCommand.DeletePermanent(ctx, s.userID)
	s.NoError(err)
	s.True(success)
}

func (s *UserRepositoryTestSuite) Test11_UpdatePassword() {
	ctx := context.Background()

	// Create a fresh user for password test
	u, err := s.repo.UserCommand.Create(ctx, &requests.CreateUserRequest{
		FirstName: "Pass", LastName: "Test", Email: "password.test@example.com",
		Password: "oldpass123", ConfirmPassword: "oldpass123",
	})
	s.Require().NoError(err)

	updatedPass, err := s.repo.UserCommand.UpdatePassword(ctx, int(u.UserID), "newpass456")
	s.NoError(err)
	s.NotNil(updatedPass)

	// Cleanup
	s.repo.UserCommand.Trash(ctx, int(u.UserID))
	s.repo.UserCommand.DeletePermanent(ctx, int(u.UserID))
}

func TestUserRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserRepositoryTestSuite))
}
