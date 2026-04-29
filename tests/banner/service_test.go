package banner_test

import (
	"context"
	"testing"

	banner_cache "github.com/MamangRust/monolith-ecommerce-grpc-banner/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type BannerServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *BannerServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.Require().NotNil(s.Obs)

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Banner dependencies
	mencache := banner_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(queries)

	s.svc = service.NewService(&service.Deps{
		Cache:         mencache,
		Repository:    repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	s.Require().NotNil(s.svc)
	s.Require().NotNil(s.svc.BannerCommand)
	s.Require().NotNil(s.svc.BannerQuery)
}

func (s *BannerServiceTestSuite) TestBannerLifecycle() {
	ctx := context.Background()

	// 1. Create
	req := &requests.CreateBannerRequest{
		Name:      "Initial Banner",
		StartDate: "2026-01-01",
		EndDate:   "2026-12-31",
		StartTime: "00:00:00",
		EndTime:   "23:59:59",
		IsActive:  true,
	}
	created, err := s.svc.BannerCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	bannerID := int(created.BannerID)

	// 2. FindByID
	found, err := s.svc.BannerQuery.FindByID(ctx, bannerID)
	s.Require().NoError(err)
	s.Equal(req.Name, found.Name)

	// 3. Update
	updateReq := &requests.UpdateBannerRequest{
		BannerID:  &bannerID,
		Name:      "Updated Banner Name",
		StartDate: "2026-01-02",
		EndDate:   "2026-12-30",
		StartTime: "01:00:00",
		EndTime:   "22:00:00",
		IsActive:  true,
	}
	updated, err := s.svc.BannerCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(updateReq.Name, updated.Name)

	// 4. FindAll
	_, total, err := s.svc.BannerQuery.FindAll(ctx, &requests.FindAllBanner{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.svc.BannerCommand.Trash(ctx, bannerID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.svc.BannerQuery.FindTrashed(ctx, &requests.FindAllBanner{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive (should be empty if only one banner was created and now trashed)
	active, _, err := s.svc.BannerQuery.FindActive(ctx, &requests.FindAllBanner{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	// Other banners might exist from other tests if not cleaned up, but our banner should not be here
	for _, b := range active {
		s.NotEqual(bannerID, int(b.BannerID))
	}

	// 8. Restore
	_, err = s.svc.BannerCommand.Restore(ctx, bannerID)
	s.Require().NoError(err)

	// 9. DeletePermanent (Trash first)
	_, err = s.svc.BannerCommand.Trash(ctx, bannerID)
	s.Require().NoError(err)
	success, err := s.svc.BannerCommand.DeletePermanent(ctx, bannerID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	// Create multiple
	b1, _ := s.svc.BannerCommand.Create(ctx, &requests.CreateBannerRequest{Name: "B1", StartDate: "2026-01-01", EndDate: "2026-12-31", IsActive: true, StartTime: "00:00:00", EndTime: "23:59:59"})
	b2, _ := s.svc.BannerCommand.Create(ctx, &requests.CreateBannerRequest{Name: "B2", StartDate: "2026-01-01", EndDate: "2026-12-31", IsActive: true, StartTime: "00:00:00", EndTime: "23:59:59"})
	
	s.svc.BannerCommand.Trash(ctx, int(b1.BannerID))
	s.svc.BannerCommand.Trash(ctx, int(b2.BannerID))

	resRestoreAll, err := s.svc.BannerCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.BannerCommand.Trash(ctx, int(b1.BannerID))
	s.svc.BannerCommand.Trash(ctx, int(b2.BannerID))

	resDeleteAll, err := s.svc.BannerCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestBannerServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(BannerServiceTestSuite))
}
