package merchant_test

import (
	"context"
	"testing"

	merchant_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantServiceTestSuite struct {
	tests.BaseTestSuite
	merchantService service.Service
	userID          int
	merchantID      int
}

func (s *MerchantServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	s.SetupUserService()
	queries := db.New(s.DBPool())
	repos := repository.NewRepositories(queries, pb.NewUserQueryServiceClient(s.Conns["user"]))

	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)

	mencache := merchant_cache.NewMencache(cacheStore)

	s.merchantService = *service.NewService(&service.Deps{
		Kafka:         nil,
		Repositories:  repos,
		Logger:        s.Log,
		Mencache:      mencache,
		Observability: s.Obs,
	})

	// 1. Seed dependencies
	s.userID = s.SeedUser(context.Background())
}

func (s *MerchantServiceTestSuite) TestMerchantLifecycle() {
	ctx := context.Background()

	req := &requests.CreateMerchantRequest{
		UserID:       s.userID,
		Name:         "Service Merchant",
		Description:  "Detailed description of the merchant.",
		Address:      "Service Street No. 1",
		ContactEmail: "service.merchant@example.com",
		ContactPhone: "08123456781",
		Status:       "active",
	}
	merchant, err := s.merchantService.MerchantCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(merchant)
	s.Equal(req.Name, merchant.Name)
	s.merchantID = int(merchant.MerchantID)

	found, err := s.merchantService.MerchantQuery.FindByID(ctx, s.merchantID)
	s.NoError(err)
	s.Equal(s.merchantID, int(found.MerchantID))
}

func TestMerchantServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantServiceTestSuite))
}
