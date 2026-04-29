package banner_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/utils"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type BannerRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *BannerRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(queries)
}

func (s *BannerRepositoryTestSuite) TestBannerLifecycle() {
	ctx := context.Background()

	req := &requests.CreateBannerRequest{
		Name:      "Summer Sale",
		StartDate: "2026-01-01",
		EndDate:   "2026-12-31",
		StartTime: "00:00:00",
		EndTime:   "23:59:59",
		IsActive:  true,
	}

	created, err := s.repo.BannerCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.Name, created.Name)

	found, err := s.repo.BannerQuery.FindByID(ctx, int(created.BannerID))
	s.NoError(err)
	s.Equal(created.Name, found.Name)

	// Update
	updateReq := &requests.UpdateBannerRequest{
		BannerID:  utils.Ptr(int(created.BannerID)),
		Name:      "Summer Sale Updated",
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		IsActive:  req.IsActive,
	}
	updated, err := s.repo.BannerCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.Name, updated.Name)

	// FindAll
	banners, err := s.repo.BannerQuery.FindAll(ctx, &requests.FindAllBanner{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(banners)

	// Trash
	_, err = s.repo.BannerCommand.Trash(ctx, int(created.BannerID))
	s.NoError(err)

	// FindTrashed
	trashed, err := s.repo.BannerQuery.FindTrashed(ctx, &requests.FindAllBanner{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashed)

	// Restore
	_, err = s.repo.BannerCommand.Restore(ctx, int(created.BannerID))
	s.NoError(err)

	// FindActive
	active, err := s.repo.BannerQuery.FindActive(ctx, &requests.FindAllBanner{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(active)

	// Trash again before permanent delete
	_, err = s.repo.BannerCommand.Trash(ctx, int(created.BannerID))
	s.NoError(err)

	// Delete Permanent
	_, err = s.repo.BannerCommand.DeletePermanent(ctx, int(created.BannerID))
	s.NoError(err)

	// FindByID (should not found or return error)
	_, err = s.repo.BannerQuery.FindByID(ctx, int(created.BannerID))
	s.Error(err)
}

func TestBannerRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(BannerRepositoryTestSuite))
}
