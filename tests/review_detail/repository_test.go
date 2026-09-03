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

func (s *ReviewDetailRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *ReviewDetailRepositoryTestSuite) TestReviewDetailLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	reviewID := s.SeedReview(ctx, userID, productID)

	pageReq := &requests.FindAllReview{Search: "", Page: 1, PageSize: 10}

	// 1. Create
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
	detailID := int(created.ReviewDetailID)

	// 2. FindByID
	found, err := s.repo.ReviewDetailQuery.FindByID(ctx, detailID)
	s.NoError(err)
	s.Equal(created.Caption, found.Caption)

	// 3. FindAll
	all, err := s.repo.ReviewDetailQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 4. FindActive (before trash)
	active, err := s.repo.ReviewDetailQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 5. Update
	updateReq := &requests.UpdateReviewDetailRequest{
		ReviewDetailID: &detailID,
		Type:           "video",
		Url:            "http://example.com/vid.jpg",
		Caption:        "Updated caption",
	}
	updated, err := s.repo.ReviewDetailCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated caption", *updated.Caption)

	// 6. Trash
	trashed, err := s.repo.ReviewDetailCommand.Trash(ctx, detailID)
	s.NoError(err)
	s.NotNil(trashed)

	// 7. FindTrashed
	trashedItems, err := s.repo.ReviewDetailQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 8. FindActive after trash — should NOT include
	activeAfterTrash, err := s.repo.ReviewDetailQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(detailID, int(item.ReviewDetailID))
	}

	// 9. Restore
	restored, err := s.repo.ReviewDetailCommand.Restore(ctx, detailID)
	s.NoError(err)
	s.NotNil(restored)

	// 10. Trash again then DeletePermanent
	_, err = s.repo.ReviewDetailCommand.Trash(ctx, detailID)
	s.Require().NoError(err)

	success, err := s.repo.ReviewDetailCommand.DeletePermanent(ctx, detailID)
	s.NoError(err)
	s.True(success)

	// 11. RestoreAll
	s.repo.ReviewDetailCommand.Create(ctx, &requests.CreateReviewDetailRequest{
		ReviewID: reviewID, Type: "photo", Url: "http://example.com/a.jpg", Caption: "A",
	})
	second, _ := s.repo.ReviewDetailCommand.Create(ctx, &requests.CreateReviewDetailRequest{
		ReviewID: reviewID, Type: "photo", Url: "http://example.com/b.jpg", Caption: "B",
	})
	s.repo.ReviewDetailCommand.Trash(ctx, int(second.ReviewDetailID))

	resRestore, err := s.repo.ReviewDetailCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 12. DeleteAll
	s.repo.ReviewDetailCommand.Trash(ctx, int(second.ReviewDetailID))
	resDelete, err := s.repo.ReviewDetailCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestReviewDetailRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewDetailRepositoryTestSuite))
}
