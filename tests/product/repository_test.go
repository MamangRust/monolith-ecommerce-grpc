package product_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ProductRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *ProductRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(
		queries,
		pb.NewCategoryQueryServiceClient(s.Conns["category"]),
		pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
}

func (s *ProductRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *ProductRepositoryTestSuite) TestProductLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	rating := 5
	slug := "test-product"
	pageReq := &requests.FindAllProduct{Search: "", Page: 1, PageSize: 10}

	// 1. Create
	req := &requests.CreateProductRequest{
		Name:         "Test Product",
		Description:  "Product for repository test",
		Price:        100,
		CountInStock: 50,
		CategoryID:   catID,
		MerchantID:   merchID,
		Brand:        "Test Brand",
		Weight:       1,
		Rating:       &rating,
		SlugProduct:  &slug,
		ImageProduct: "test.jpg",
	}
	created, err := s.repo.ProductCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.Name, created.Name)
	prodID := int(created.ProductID)

	// 2. FindByID
	found, err := s.repo.ProductQuery.FindByID(ctx, prodID)
	s.NoError(err)
	s.Equal(created.Name, found.Name)

	// 3. FindAll
	all, err := s.repo.ProductQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 4. FindActive (before trash)
	active, err := s.repo.ProductQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 5. Update
	updatedRating := 4
	updatedSlug := "updated-product"
	updateReq := &requests.UpdateProductRequest{
		ProductID:    &prodID,
		MerchantID:   merchID,
		CategoryID:   catID,
		Name:         "Updated Product",
		Description:  "Updated description",
		Price:        200,
		CountInStock: 100,
		Brand:        "Updated Brand",
		Weight:       2,
		Rating:       &updatedRating,
		SlugProduct:  &updatedSlug,
		ImageProduct: "updated.jpg",
	}
	updated, err := s.repo.ProductCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated Product", updated.Name)

	// 6. Trash
	trashed, err := s.repo.ProductCommand.Trash(ctx, prodID)
	s.NoError(err)
	s.NotNil(trashed)

	// 7. FindTrashed
	trashedItems, err := s.repo.ProductQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 8. FindActive after trash — should NOT include
	activeAfterTrash, err := s.repo.ProductQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(prodID, int(item.ProductID))
	}

	// 9. Restore
	restored, err := s.repo.ProductCommand.Restore(ctx, prodID)
	s.NoError(err)
	s.NotNil(restored)

	// 10. Trash again then DeletePermanent
	_, err = s.repo.ProductCommand.Trash(ctx, prodID)
	s.Require().NoError(err)

	success, err := s.repo.ProductCommand.DeletePermanent(ctx, prodID)
	s.NoError(err)
	s.True(success)

	// 11. RestoreAll
	slug1 := "p1"
	slug2 := "p2"
	r1 := 5
	r2 := 4
	s.repo.ProductCommand.Create(ctx, &requests.CreateProductRequest{
		Name: "P1", Description: "D1", Price: 100, CountInStock: 10,
		CategoryID: catID, MerchantID: merchID, Brand: "B", Weight: 1,
		Rating: &r1, SlugProduct: &slug1, ImageProduct: "img.jpg",
	})
	second, _ := s.repo.ProductCommand.Create(ctx, &requests.CreateProductRequest{
		Name: "P2", Description: "D2", Price: 200, CountInStock: 20,
		CategoryID: catID, MerchantID: merchID, Brand: "B", Weight: 2,
		Rating: &r2, SlugProduct: &slug2, ImageProduct: "img.jpg",
	})
	s.repo.ProductCommand.Trash(ctx, int(second.ProductID))

	resRestore, err := s.repo.ProductCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 12. DeleteAll
	s.repo.ProductCommand.Trash(ctx, int(second.ProductID))
	resDelete, err := s.repo.ProductCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func (s *ProductRepositoryTestSuite) TestProductFindByIDNotFound() {
	ctx := context.Background()

	// FindByID on a non-existent ID must return a typed not-found error.
	_, err := s.repo.ProductQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}

func TestProductRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductRepositoryTestSuite))
}
