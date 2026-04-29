package category_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/utils"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CategoryRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *CategoryRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(queries)
}

func (s *CategoryRepositoryTestSuite) TestCategoryLifecycle() {
	ctx := context.Background()
	req := &requests.CreateCategoryRequest{
		Name:          "Electronics",
		Description:   "Electronic items and gadgets",
		ImageCategory: "electronics.jpg",
	}

	// Create
	created, err := s.repo.CategoryCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.Name, created.Name)

	// Find by ID
	found, err := s.repo.CategoryQuery.FindByID(ctx, int(created.CategoryID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(created.Name, found.Name)

	// Update
	updateReq := &requests.UpdateCategoryRequest{
		CategoryID:    utils.Ptr(int(created.CategoryID)),
		Name:          "Electronics Updated",
		Description:   req.Description,
		ImageCategory: req.ImageCategory,
	}
	updated, err := s.repo.CategoryCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.Name, updated.Name)

	// Find All
	all, err := s.repo.CategoryQuery.FindAll(ctx, &requests.FindAllCategory{PageSize: 10, Page: 1})
	s.NoError(err)
	s.NotEmpty(all)

	// Trash
	_, err = s.repo.CategoryCommand.Trash(ctx, int(created.CategoryID))
	s.NoError(err)

	// Find Trashed
	trashed, err := s.repo.CategoryQuery.FindTrashed(ctx, &requests.FindAllCategory{PageSize: 10, Page: 1})
	s.NoError(err)
	s.NotEmpty(trashed)

	// Restore
	_, err = s.repo.CategoryCommand.Restore(ctx, int(created.CategoryID))
	s.NoError(err)

	// Find Active
	active, err := s.repo.CategoryQuery.FindActive(ctx, &requests.FindAllCategory{PageSize: 10, Page: 1})
	s.NoError(err)
	s.NotEmpty(active)

	// Trash again before permanent delete
	_, err = s.repo.CategoryCommand.Trash(ctx, int(created.CategoryID))
	s.NoError(err)

	// Delete Permanent
	_, err = s.repo.CategoryCommand.DeletePermanent(ctx, int(created.CategoryID))
	s.NoError(err)

	// FindByID (should not found or return error)
	_, err = s.repo.CategoryQuery.FindByID(ctx, int(created.CategoryID))
	s.Error(err)
}

func TestCategoryRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryRepositoryTestSuite))
}
