package merchant_policy_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantPolicyRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *MerchantPolicyRepositoryTestSuite) SetupSuite() {
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

func (s *MerchantPolicyRepositoryTestSuite) TestPolicyLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	// 2. Create Policy
	req := &requests.CreateMerchantPolicyRequest{
		MerchantID:  merchID,
		PolicyType:  "Return",
		Title:       "Return Policy",
		Description: "No returns after 7 days",
	}

	created, err := s.repo.MerchantPoliciesCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.PolicyType, created.PolicyType)

	// 3. Find by ID
	found, err := s.repo.MerchantPoliciesQuery.FindByID(ctx, int(created.MerchantPolicyID))
	s.NoError(err)
	s.Equal(created.MerchantPolicyID, found.MerchantPolicyID)
}

func TestMerchantPolicyRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantPolicyRepositoryTestSuite))
}
