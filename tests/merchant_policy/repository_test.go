package merchant_policy_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
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

func (s *MerchantPolicyRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *MerchantPolicyRepositoryTestSuite) TestPolicyLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	pageReq := &requests.FindAllMerchant{Search: "", Page: 1, PageSize: 10}

	// 1. Create Policy
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
	policyID := int(created.MerchantPolicyID)

	// 2. FindByID
	found, err := s.repo.MerchantPoliciesQuery.FindByID(ctx, policyID)
	s.NoError(err)
	s.Equal(created.MerchantPolicyID, found.MerchantPolicyID)

	// 3. FindAll
	all, err := s.repo.MerchantPoliciesQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 4. FindActive (before trash)
	active, err := s.repo.MerchantPoliciesQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 5. Update
	updateReq := &requests.UpdateMerchantPolicyRequest{
		MerchantPolicyID: &policyID,
		PolicyType:       "Warranty",
		Title:            "Updated Policy",
		Description:      "Updated description",
	}
	updated, err := s.repo.MerchantPoliciesCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Warranty", updated.PolicyType)

	// 6. Trash
	trashed, err := s.repo.MerchantPoliciesCommand.Trash(ctx, policyID)
	s.NoError(err)
	s.NotNil(trashed)

	// 7. FindTrashed
	trashedItems, err := s.repo.MerchantPoliciesQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 8. FindActive after trash — should NOT include
	activeAfterTrash, err := s.repo.MerchantPoliciesQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(policyID, int(item.MerchantPolicyID))
	}

	// 9. Restore
	restored, err := s.repo.MerchantPoliciesCommand.Restore(ctx, policyID)
	s.NoError(err)
	s.NotNil(restored)

	// 10. Trash again then DeletePermanent
	_, err = s.repo.MerchantPoliciesCommand.Trash(ctx, policyID)
	s.Require().NoError(err)

	success, err := s.repo.MerchantPoliciesCommand.DeletePermanent(ctx, policyID)
	s.NoError(err)
	s.True(success)

	// 11. RestoreAll
	s.repo.MerchantPoliciesCommand.Create(ctx, &requests.CreateMerchantPolicyRequest{
		MerchantID: merchID, PolicyType: "T1", Title: "T1", Description: "D1",
	})
	second, _ := s.repo.MerchantPoliciesCommand.Create(ctx, &requests.CreateMerchantPolicyRequest{
		MerchantID: merchID, PolicyType: "T2", Title: "T2", Description: "D2",
	})
	s.repo.MerchantPoliciesCommand.Trash(ctx, int(second.MerchantPolicyID))

	resRestore, err := s.repo.MerchantPoliciesCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 12. DeleteAll
	s.repo.MerchantPoliciesCommand.Trash(ctx, int(second.MerchantPolicyID))
	resDelete, err := s.repo.MerchantPoliciesCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestMerchantPolicyRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantPolicyRepositoryTestSuite))
}
