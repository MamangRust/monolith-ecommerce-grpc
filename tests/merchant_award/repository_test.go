package merchant_award_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantAwardRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *MerchantAwardRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup core merchant service dependency
	s.SetupUserService()
	s.SetupMerchantService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(
		queries,
		pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
}

func (s *MerchantAwardRepositoryTestSuite) TestAwardLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	// 2. Create Award
	req := &requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID:  merchID,
		Title:       "Best Merchant 2024",
		Description: "Detailed description of the achievement.",
		IssueDate:   "2024-01-01",
		IssuedBy:    "Ecommerce Platform",
	}

	created, err := s.repo.MerchantAwardCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.Title, created.Title)

	// 3. Find by ID
	found, err := s.repo.MerchantAwardQuery.FindByID(ctx, int(created.MerchantCertificationID))
	s.NoError(err)
	s.Equal(created.Title, found.Title)
}

func TestMerchantAwardRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantAwardRepositoryTestSuite))
}
