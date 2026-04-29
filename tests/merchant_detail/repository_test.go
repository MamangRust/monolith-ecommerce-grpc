package merchant_detail_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantDetailRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *MerchantDetailRepositoryTestSuite) SetupSuite() {
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

func (s *MerchantDetailRepositoryTestSuite) TestDetailLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	// 2. Create Detail
	req := &requests.CreateMerchantDetailRequest{
		MerchantID:      merchID,
		DisplayName:     "Detail Display",
		ShortDescription: "Detail Description",
	}

	created, err := s.repo.MerchantDetailCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.DisplayName, *created.DisplayName)

	// 3. Find by ID
	found, err := s.repo.MerchantDetailQuery.FindByID(ctx, int(created.MerchantDetailID))
	s.NoError(err)
	s.Equal(created.MerchantDetailID, found.MerchantDetailID)
}

func TestMerchantDetailRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDetailRepositoryTestSuite))
}
