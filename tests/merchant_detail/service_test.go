package merchant_detail_test

import (
	"context"
	"testing"

	detail_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantDetailServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *MerchantDetailServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.Require().NotNil(s.Obs)

	// Setup core dependencies
	s.SetupUserService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Detail dependencies
	mencache := detail_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(
		queries,
		pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)

	s.svc = service.NewService(&service.Deps{
		Cache:         mencache,
		Repository:    repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	s.Require().NotNil(s.svc)
	s.Require().NotNil(s.svc.MerchantDetailCommand)
	s.Require().NotNil(s.svc.MerchantDetailQuery)
}

func (s *MerchantDetailServiceTestSuite) TestMerchantDetailLifecycle() {
	ctx := context.Background()

	// Setup Merchant
	s.SetupUserService()
	s.SetupMerchantService()
	userID := s.SeedUser(ctx)
	merchantID := s.SeedMerchant(ctx, userID)

	// 1. Create
	req := &requests.CreateMerchantDetailRequest{
		MerchantID:       int(merchantID),
		DisplayName:      "Test Display Name",
		ShortDescription: "Test Short Description",
		CoverImageUrl:    "http://example.com/cover.jpg",
		LogoUrl:          "http://example.com/logo.jpg",
		WebsiteUrl:       "http://example.com",
	}
	created, err := s.svc.MerchantDetailCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	detailID := int(created.MerchantDetailID)

	// 2. FindByID
	found, err := s.svc.MerchantDetailQuery.FindByID(ctx, int(merchantID)) // Note: FindByID takes merchantID
	s.Require().NoError(err)
	s.Equal(req.DisplayName, *found.DisplayName)

	// 3. Update
	newDisplayName := "Updated Display Name"
	updateReq := &requests.UpdateMerchantDetailRequest{
		MerchantDetailID: &detailID,
		DisplayName:      newDisplayName,
		ShortDescription: req.ShortDescription,
		CoverImageUrl:    req.CoverImageUrl,
		LogoUrl:          req.LogoUrl,
		WebsiteUrl:       req.WebsiteUrl,
	}
	updated, err := s.svc.MerchantDetailCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(newDisplayName, *updated.DisplayName)

	// 4. FindAll
	_, total, err := s.svc.MerchantDetailQuery.FindAll(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.svc.MerchantDetailCommand.Trash(ctx, detailID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.svc.MerchantDetailQuery.FindTrashed(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.svc.MerchantDetailQuery.FindActive(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, d := range active {
		s.NotEqual(detailID, int(d.MerchantDetailID))
	}

	// 8. Restore
	_, err = s.svc.MerchantDetailCommand.Restore(ctx, detailID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.svc.MerchantDetailCommand.Trash(ctx, detailID)
	s.Require().NoError(err)
	success, err := s.svc.MerchantDetailCommand.DeletePermanent(ctx, detailID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	d1, _ := s.svc.MerchantDetailCommand.Create(ctx, &requests.CreateMerchantDetailRequest{MerchantID: int(merchantID), DisplayName: "D1", ShortDescription: "SD1"})
	d2, _ := s.svc.MerchantDetailCommand.Create(ctx, &requests.CreateMerchantDetailRequest{MerchantID: int(merchantID), DisplayName: "D2", ShortDescription: "SD2"})
	
	s.svc.MerchantDetailCommand.Trash(ctx, int(d1.MerchantDetailID))
	s.svc.MerchantDetailCommand.Trash(ctx, int(d2.MerchantDetailID))

	resRestoreAll, err := s.svc.MerchantDetailCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.MerchantDetailCommand.Trash(ctx, int(d1.MerchantDetailID))
	s.svc.MerchantDetailCommand.Trash(ctx, int(d2.MerchantDetailID))

	resDeleteAll, err := s.svc.MerchantDetailCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestMerchantDetailServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDetailServiceTestSuite))
}
