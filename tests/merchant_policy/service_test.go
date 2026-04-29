package merchant_policy_test

import (
	"context"
	"testing"

	policy_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantPolicyServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *MerchantPolicyServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.Require().NotNil(s.Obs)

	// Setup core dependencies
	s.SetupUserService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Policy dependencies
	mencache := policy_cache.NewMencache(cacheStore)
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

func (s *MerchantPolicyServiceTestSuite) TestMerchantPolicyLifecycle() {
	ctx := context.Background()

	// Setup Merchant
	s.SetupUserService()
	s.SetupMerchantService()
	userID := s.SeedUser(ctx)
	merchantID := s.SeedMerchant(ctx, userID)

	// 1. Create
	req := &requests.CreateMerchantPolicyRequest{
		MerchantID:  int(merchantID),
		PolicyType:  "Shipping",
		Title:       "Shipping Policy",
		Description: "Standard shipping in 3-5 days.",
	}
	created, err := s.svc.MerchantPoliciesCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	policyID := int(created.MerchantPolicyID)

	// 2. FindByID
	found, err := s.svc.MerchantPoliciesQuery.FindByID(ctx, policyID)
	s.Require().NoError(err)
	s.Equal(req.Title, found.Title)

	// 3. Update
	newTitle := "Express Shipping Policy"
	updateReq := &requests.UpdateMerchantPolicyRequest{
		MerchantPolicyID: &policyID,
		PolicyType:       req.PolicyType,
		Title:            newTitle,
		Description:      req.Description,
	}
	updated, err := s.svc.MerchantPoliciesCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(newTitle, updated.Title)

	// 4. FindAll
	_, total, err := s.svc.MerchantPoliciesQuery.FindAll(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.svc.MerchantPoliciesCommand.Trash(ctx, policyID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.svc.MerchantPoliciesQuery.FindTrashed(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.svc.MerchantPoliciesQuery.FindActive(ctx, &requests.FindAllMerchant{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, p := range active {
		s.NotEqual(policyID, int(p.MerchantPolicyID))
	}

	// 8. Restore
	_, err = s.svc.MerchantPoliciesCommand.Restore(ctx, policyID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.svc.MerchantPoliciesCommand.Trash(ctx, policyID)
	s.Require().NoError(err)
	success, err := s.svc.MerchantPoliciesCommand.DeletePermanent(ctx, policyID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	p1, _ := s.svc.MerchantPoliciesCommand.Create(ctx, &requests.CreateMerchantPolicyRequest{MerchantID: int(merchantID), PolicyType: "T1", Title: "T1", Description: "D1"})
	p2, _ := s.svc.MerchantPoliciesCommand.Create(ctx, &requests.CreateMerchantPolicyRequest{MerchantID: int(merchantID), PolicyType: "T2", Title: "T2", Description: "D2"})
	
	s.svc.MerchantPoliciesCommand.Trash(ctx, int(p1.MerchantPolicyID))
	s.svc.MerchantPoliciesCommand.Trash(ctx, int(p2.MerchantPolicyID))

	resRestoreAll, err := s.svc.MerchantPoliciesCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.MerchantPoliciesCommand.Trash(ctx, int(p1.MerchantPolicyID))
	s.svc.MerchantPoliciesCommand.Trash(ctx, int(p2.MerchantPolicyID))

	resDeleteAll, err := s.svc.MerchantPoliciesCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestMerchantPolicyServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantPolicyServiceTestSuite))
}
