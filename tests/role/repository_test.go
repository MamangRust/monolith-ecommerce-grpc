package role_test

import (
	"context"
	"testing"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type RoleRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	repo   *repository.Repositories
	roleID int
}

func (s *RoleRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)

	queries := db.New(pool)
	s.repo = repository.NewRepositories(queries)
}

func (s *RoleRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *RoleRepositoryTestSuite) Test1_CreateRole() {
	ctx := context.Background()
	req := &requests.CreateRoleRequest{
		Name: "Test Role",
	}

	res, err := s.repo.RoleCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Name, res.RoleName)
	s.roleID = int(res.RoleID)
}

func (s *RoleRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	found, err := s.repo.RoleQuery.FindByID(ctx, s.roleID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.roleID, int(found.RoleID))
	s.Equal("Test Role", found.RoleName)
}

func (s *RoleRepositoryTestSuite) Test3_UpdateRole() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	req := &requests.UpdateRoleRequest{
		ID:   &s.roleID,
		Name: "Updated Role",
	}

	res, err := s.repo.RoleCommand.Update(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("Updated Role", res.RoleName)
}

func (s *RoleRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	// Trash
	trashed, err := s.repo.RoleCommand.Trash(ctx, s.roleID)
	s.NoError(err)
	s.NotNil(trashed)

	// Restore
	restored, err := s.repo.RoleCommand.Restore(ctx, s.roleID)
	s.NoError(err)
	s.NotNil(restored)

	// Verify restored
	found, err := s.repo.RoleQuery.FindByID(ctx, s.roleID)
	s.NoError(err)
	s.NotNil(found)
}

func (s *RoleRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	// Must be trashed first for permanent delete
	_, err := s.repo.RoleCommand.Trash(ctx, s.roleID)
	s.NoError(err)

	success, err := s.repo.RoleCommand.DeletePermanent(ctx, s.roleID)
	s.NoError(err)
	s.True(success)
}

func (s *RoleRepositoryTestSuite) Test6_TrashMissingRole_ReturnsNotFound() {
	ctx := context.Background()

	// Trashing a role that does not exist must surface as ErrRoleNotFound
	// (HTTP 404), not as a generic internal error (HTTP 500).
	res, err := s.repo.RoleCommand.Trash(ctx, 999999)
	s.Nil(res)
	s.Require().Error(err)
	s.ErrorIs(err, role_errors.ErrRoleNotFound)
}

func (s *RoleRepositoryTestSuite) Test7_RestoreMissingRole_ReturnsNotFound() {
	ctx := context.Background()

	// Restoring a role that does not exist must surface as ErrRoleNotFound
	// (HTTP 404), not as a generic internal error (HTTP 500).
	res, err := s.repo.RoleCommand.Restore(ctx, 999999)
	s.Nil(res)
	s.Require().Error(err)
	s.ErrorIs(err, role_errors.ErrRoleNotFound)
}

func (s *RoleRepositoryTestSuite) Test8_CreateDuplicateRole_ReturnsConflict() {
	ctx := context.Background()

	// Creating a role with a name that already exists must surface as
	// ErrRoleConflict (HTTP 409), not as a generic internal error (HTTP 500).
	name := "Duplicate Role"
	_, err := s.repo.RoleCommand.Create(ctx, &requests.CreateRoleRequest{Name: name})
	s.Require().NoError(err)

	res, err := s.repo.RoleCommand.Create(ctx, &requests.CreateRoleRequest{Name: name})
	s.Nil(res)
	s.Require().Error(err)
	s.ErrorIs(err, role_errors.ErrRoleConflict)
}

func (s *RoleRepositoryTestSuite) Test9_UpdateDuplicateRole_ReturnsConflict() {
	ctx := context.Background()

	// Renaming a role to a name that already exists must surface as
	// ErrRoleConflict (HTTP 409), not as a generic internal error (HTTP 500).
	original, err := s.repo.RoleCommand.Create(ctx, &requests.CreateRoleRequest{Name: "Original Role"})
	s.Require().NoError(err)

	other, err := s.repo.RoleCommand.Create(ctx, &requests.CreateRoleRequest{Name: "Other Role"})
	s.Require().NoError(err)

	otherID := int(other.RoleID)
	res, err := s.repo.RoleCommand.Update(ctx, &requests.UpdateRoleRequest{
		ID:   &otherID,
		Name: original.RoleName,
	})
	s.Nil(res)
	s.Require().Error(err)
	s.ErrorIs(err, role_errors.ErrRoleConflict)
}

func TestRoleRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleRepositoryTestSuite))
}
