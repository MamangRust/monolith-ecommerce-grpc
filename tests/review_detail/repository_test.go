package review_detail_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ReviewDetailRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *ReviewDetailRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupOrderService()
	s.SetupReviewService()
	s.SetupReviewDetailService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(queries)
}

func (s *ReviewDetailRepositoryTestSuite) TestReviewDetailLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	reviewID := s.SeedReview(ctx, userID, productID)

	// 2. Create Review Detail
	req := &requests.CreateReviewDetailRequest{
		ReviewID: reviewID,
		Type:     "photo",
		Url:      "http://example.com/img.jpg",
		Caption:  "This is a detailed feedback with more info.",
	}

	created, err := s.repo.ReviewDetailCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(reviewID), created.ReviewID)

	// 3. Find by ID
	found, err := s.repo.ReviewDetailQuery.FindByID(ctx, int(created.ReviewDetailID))
	s.NoError(err)
	s.Equal(created.Caption, found.Caption)

	// 4. Find by Review ID
	details, err := s.repo.ReviewDetailQuery.FindByID(ctx, reviewID)
	s.NoError(err)
	s.NotEmpty(details)
}

func TestReviewDetailRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewDetailRepositoryTestSuite))
}
