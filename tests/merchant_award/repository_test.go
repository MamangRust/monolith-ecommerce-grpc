package merchant_award_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantAwardRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *MerchantAwardRepositoryTestSuite) SetupSuite() {
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

func (s *MerchantAwardRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *MerchantAwardRepositoryTestSuite) TestAwardLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	pageReq := &requests.FindAllMerchant{Search: "", Page: 1, PageSize: 10}

	// 1. Create Award
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
	awardID := int(created.MerchantCertificationID)

	// 2. FindByID
	found, err := s.repo.MerchantAwardQuery.FindByID(ctx, awardID)
	s.NoError(err)
	s.Equal(created.Title, found.Title)

	// 3. FindAll
	all, err := s.repo.MerchantAwardQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 4. FindActive (before trash)
	active, err := s.repo.MerchantAwardQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 5. Update
	updateReq := &requests.UpdateMerchantCertificationOrAwardRequest{
		MerchantCertificationID: &awardID,
		Title:                   "Updated Award",
		Description:             "Updated description",
		IssuedBy:                "Updated Platform",
		IssueDate:               "2024-06-01",
	}
	updated, err := s.repo.MerchantAwardCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated Award", updated.Title)

	// 6. Trash
	trashed, err := s.repo.MerchantAwardCommand.Trash(ctx, awardID)
	s.NoError(err)
	s.NotNil(trashed)

	// 7. FindTrashed
	trashedItems, err := s.repo.MerchantAwardQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 8. FindActive after trash — should NOT include
	activeAfterTrash, err := s.repo.MerchantAwardQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(awardID, int(item.MerchantCertificationID))
	}

	// 9. Restore
	restored, err := s.repo.MerchantAwardCommand.Restore(ctx, awardID)
	s.NoError(err)
	s.NotNil(restored)

	// 10. Trash again then DeletePermanent
	_, err = s.repo.MerchantAwardCommand.Trash(ctx, awardID)
	s.Require().NoError(err)

	success, err := s.repo.MerchantAwardCommand.DeletePermanent(ctx, awardID)
	s.NoError(err)
	s.True(success)

	// 11. RestoreAll
	s.repo.MerchantAwardCommand.Create(ctx, &requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID: merchID, Title: "A1", Description: "D1", IssueDate: "2024-01-01", IssuedBy: "P1",
	})
	second, _ := s.repo.MerchantAwardCommand.Create(ctx, &requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID: merchID, Title: "A2", Description: "D2", IssueDate: "2024-02-01", IssuedBy: "P2",
	})
	s.repo.MerchantAwardCommand.Trash(ctx, int(second.MerchantCertificationID))

	resRestore, err := s.repo.MerchantAwardCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 12. DeleteAll
	s.repo.MerchantAwardCommand.Trash(ctx, int(second.MerchantCertificationID))
	resDelete, err := s.repo.MerchantAwardCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestMerchantAwardRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantAwardRepositoryTestSuite))
}
