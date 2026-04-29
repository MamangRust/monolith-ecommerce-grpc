package review_detail_test

import (
	"context"
	"testing"

	detail_cache "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ReviewDetailServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *ReviewDetailServiceTestSuite) SetupSuite() {
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

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Review Detail dependencies
	mencache := detail_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(queries)

	s.svc = service.NewService(&service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
}

func (s *ReviewDetailServiceTestSuite) TestReviewDetailLifecycle() {
	ctx := context.Background()

	// 1. Setup Dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupReviewService()

	userID := s.SeedUser(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	categoryID := s.SeedCategory(ctx)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	reviewID := s.SeedReview(ctx, userID, productID)

	// 2. Create Review Detail
	req := &requests.CreateReviewDetailRequest{
		ReviewID: reviewID,
		Type:     "photo",
		Url:      "detail_photo.jpg",
		Caption:  "Initial Caption",
	}
	created, err := s.svc.ReviewDetailCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	detailID := int(created.ReviewDetailID)

	// 3. FindByID
	found, err := s.svc.ReviewDetailQuery.FindByID(ctx, detailID)
	s.Require().NoError(err)
	s.Equal(req.Url, found.Url)

	// 4. Update
	newUrl := "updated_detail_photo.jpg"
	updateReq := &requests.UpdateReviewDetailRequest{
		ReviewDetailID: &detailID,
		Type:           req.Type,
		Url:            newUrl,
		Caption:        req.Caption,
	}
	updated, err := s.svc.ReviewDetailCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(newUrl, updated.Url)

	// 5. FindAll
	_, total, err := s.svc.ReviewDetailQuery.FindAll(ctx, &requests.FindAllReview{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 6. Trash
	_, err = s.svc.ReviewDetailCommand.Trash(ctx, detailID)
	s.Require().NoError(err)

	// 7. FindTrashed
	_, totalTrashed, err := s.svc.ReviewDetailQuery.FindTrashed(ctx, &requests.FindAllReview{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 8. FindActive
	active, _, err := s.svc.ReviewDetailQuery.FindActive(ctx, &requests.FindAllReview{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, d := range active {
		s.NotEqual(detailID, int(d.ReviewDetailID))
	}

	// 9. Restore
	_, err = s.svc.ReviewDetailCommand.Restore(ctx, detailID)
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, err = s.svc.ReviewDetailCommand.Trash(ctx, detailID)
	s.Require().NoError(err)
	success, err := s.svc.ReviewDetailCommand.DeletePermanent(ctx, detailID)
	s.Require().NoError(err)
	s.True(success)

	// 11. RestoreAll & DeleteAll
	d1, _ := s.svc.ReviewDetailCommand.Create(ctx, &requests.CreateReviewDetailRequest{ReviewID: reviewID, Type: "photo", Url: "D1", Caption: "C1"})
	d2, _ := s.svc.ReviewDetailCommand.Create(ctx, &requests.CreateReviewDetailRequest{ReviewID: reviewID, Type: "video", Url: "D2", Caption: "C2"})
	
	s.svc.ReviewDetailCommand.Trash(ctx, int(d1.ReviewDetailID))
	s.svc.ReviewDetailCommand.Trash(ctx, int(d2.ReviewDetailID))

	resRestoreAll, err := s.svc.ReviewDetailCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.ReviewDetailCommand.Trash(ctx, int(d1.ReviewDetailID))
	s.svc.ReviewDetailCommand.Trash(ctx, int(d2.ReviewDetailID))

	resDeleteAll, err := s.svc.ReviewDetailCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestReviewDetailServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewDetailServiceTestSuite))
}
