package merchant_business_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantBusinessRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *MerchantBusinessRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup dependencies
	s.SetupUserService()
	s.SetupMerchantService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(
		queries,
		pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
}

func (s *MerchantBusinessRepositoryTestSuite) TestBusinessLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	// 2. Create Business Profile
	req := &requests.CreateMerchantBusinessInformationRequest{
		MerchantID:        merchID,
		BusinessType:      "Retail",
		TaxID:             "123-TAX",
		EstablishedYear:   2021,
		NumberOfEmployees: 10,
		WebsiteUrl:        "http://mamang.corp",
	}

	created, err := s.repo.MerchantBusinessCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.BusinessType, *created.BusinessType)

	// 3. Find by ID
	found, err := s.repo.MerchantBusinessQuery.FindByID(ctx, int(created.MerchantBusinessInfoID))
	s.NoError(err)
	s.Equal(created.MerchantBusinessInfoID, found.MerchantBusinessInfoID)
}

func TestMerchantBusinessRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantBusinessRepositoryTestSuite))
}
