package merchant_award_test

import (
	"context"
	"testing"

	award_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantAwardServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *MerchantAwardServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup core dependencies
	s.SetupUserService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Award dependencies
	mencache := award_cache.NewMencache(cacheStore)
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
}

func (s *MerchantAwardServiceTestSuite) TestMerchantAwardLifecycle() {
	ctx := context.Background()

	// Setup Merchant
	s.SetupUserService()
	s.SetupMerchantService()
	userID := s.SeedUser(ctx)
	merchantID := s.SeedMerchant(ctx, userID)

	// 1. Create
	req := &requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID:     int(merchantID),
		Title:          "Initial Award",
		Description:    "Initial Description",
		IssuedBy:       "Test Issuer",
		IssueDate:      "2026-0" + "1-02",
	}
	created, err := s.svc.MerchantAwardCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	awardID := int(created.MerchantCertificationID)

	// 2. FindByID
	found, err := s.svc.MerchantAwardQuery.FindByID(ctx, awardID)
	s.Require().NoError(err)
	s.Equal(req.Title, found.Title)

	// 3. Update
	newTitle := "Updated Award Title"
	updateReq := &requests.UpdateMerchantCertificationOrAwardRequest{
		MerchantCertificationID: &awardID,
		Title:                   newTitle,
		Description:             req.Description,
		IssuedBy:                req.IssuedBy,
		IssueDate:               req.IssueDate,
	}
	updated, err := s.svc.MerchantAwardCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(newTitle, updated.Title)

	// 4. FindAll
	_, total, err := s.svc.MerchantAwardQuery.FindAll(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.svc.MerchantAwardCommand.Trash(ctx, awardID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.svc.MerchantAwardQuery.FindTrashed(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.svc.MerchantAwardQuery.FindActive(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, a := range active {
		s.NotEqual(awardID, int(a.MerchantCertificationID))
	}

	// 8. Restore
	_, err = s.svc.MerchantAwardCommand.Restore(ctx, awardID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.svc.MerchantAwardCommand.Trash(ctx, awardID)
	s.Require().NoError(err)
	success, err := s.svc.MerchantAwardCommand.DeletePermanent(ctx, awardID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	a1, _ := s.svc.MerchantAwardCommand.Create(ctx, &requests.CreateMerchantCertificationOrAwardRequest{MerchantID: int(merchantID), Title: "A1", Description: "D1", IssuedBy: "I1", IssueDate: "2026-01-01"})
	a2, _ := s.svc.MerchantAwardCommand.Create(ctx, &requests.CreateMerchantCertificationOrAwardRequest{MerchantID: int(merchantID), Title: "A2", Description: "D2", IssuedBy: "I2", IssueDate: "2026-01-01"})
	
	s.svc.MerchantAwardCommand.Trash(ctx, int(a1.MerchantCertificationID))
	s.svc.MerchantAwardCommand.Trash(ctx, int(a2.MerchantCertificationID))

	resRestoreAll, err := s.svc.MerchantAwardCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.MerchantAwardCommand.Trash(ctx, int(a1.MerchantCertificationID))
	s.svc.MerchantAwardCommand.Trash(ctx, int(a2.MerchantCertificationID))

	resDeleteAll, err := s.svc.MerchantAwardCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestMerchantAwardServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantAwardServiceTestSuite))
}
