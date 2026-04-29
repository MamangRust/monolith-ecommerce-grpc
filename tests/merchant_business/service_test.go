package merchant_business_test

import (
	"context"
	"testing"

	biz_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantBusinessServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *MerchantBusinessServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.Require().NotNil(s.Obs)

	// Setup core dependencies
	s.SetupUserService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Business dependencies
	mencache := biz_cache.NewMencache(cacheStore)
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
	s.Require().NotNil(s.svc.MerchantBusinessCommand)
	s.Require().NotNil(s.svc.MerchantBusinessQuery)
}

func (s *MerchantBusinessServiceTestSuite) TestMerchantBusinessLifecycle() {
	ctx := context.Background()

	// Setup Merchant
	s.SetupUserService()
	s.SetupMerchantService()
	userID := s.SeedUser(ctx)
	merchantID := s.SeedMerchant(ctx, userID)

	// 1. Create
	req := &requests.CreateMerchantBusinessInformationRequest{
		MerchantID:        int(merchantID),
		BusinessType:      "PT",
		TaxID:             "12.345.678.9-012.000",
		EstablishedYear:   2020,
		NumberOfEmployees: 100,
		WebsiteUrl:        "http://example.com",
	}
	created, err := s.svc.MerchantBusinessCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	businessID := int(created.MerchantBusinessInfoID)

	// 2. FindByID
	found, err := s.svc.MerchantBusinessQuery.FindByID(ctx, int(merchantID)) // Note: FindByID takes merchantID
	s.Require().NoError(err)
	s.Equal(req.BusinessType, *found.BusinessType)

	// 3. Update
	newTaxID := "99.888.777.6-543.211"
	updateReq := &requests.UpdateMerchantBusinessInformationRequest{
		MerchantBusinessInfoID: &businessID,
		BusinessType:           req.BusinessType,
		TaxID:                  newTaxID,
		EstablishedYear:        req.EstablishedYear,
		NumberOfEmployees:      req.NumberOfEmployees,
		WebsiteUrl:             req.WebsiteUrl,
	}
	updated, err := s.svc.MerchantBusinessCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(newTaxID, *updated.TaxID)

	// 4. FindAll
	_, total, err := s.svc.MerchantBusinessQuery.FindAll(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.svc.MerchantBusinessCommand.Trash(ctx, businessID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.svc.MerchantBusinessQuery.FindTrashed(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.svc.MerchantBusinessQuery.FindActive(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, b := range active {
		s.NotEqual(businessID, int(b.MerchantBusinessInfoID))
	}

	// 8. Restore
	_, err = s.svc.MerchantBusinessCommand.Restore(ctx, businessID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.svc.MerchantBusinessCommand.Trash(ctx, businessID)
	s.Require().NoError(err)
	success, err := s.svc.MerchantBusinessCommand.DeletePermanent(ctx, businessID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	b1, _ := s.svc.MerchantBusinessCommand.Create(ctx, &requests.CreateMerchantBusinessInformationRequest{MerchantID: int(merchantID), BusinessType: "PT1", TaxID: "T1", EstablishedYear: 2021, NumberOfEmployees: 10})
	b2, _ := s.svc.MerchantBusinessCommand.Create(ctx, &requests.CreateMerchantBusinessInformationRequest{MerchantID: int(merchantID), BusinessType: "PT2", TaxID: "T2", EstablishedYear: 2022, NumberOfEmployees: 20})
	
	s.svc.MerchantBusinessCommand.Trash(ctx, int(b1.MerchantBusinessInfoID))
	s.svc.MerchantBusinessCommand.Trash(ctx, int(b2.MerchantBusinessInfoID))

	resRestoreAll, err := s.svc.MerchantBusinessCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.MerchantBusinessCommand.Trash(ctx, int(b1.MerchantBusinessInfoID))
	s.svc.MerchantBusinessCommand.Trash(ctx, int(b2.MerchantBusinessInfoID))

	resDeleteAll, err := s.svc.MerchantBusinessCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestMerchantBusinessServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantBusinessServiceTestSuite))
}
