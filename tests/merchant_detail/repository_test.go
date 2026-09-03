package merchant_detail_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
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

func (s *MerchantDetailRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *MerchantDetailRepositoryTestSuite) TestDetailLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	pageReq := &requests.FindAllMerchant{Search: "", Page: 1, PageSize: 10}

	// 1. Create Detail
	req := &requests.CreateMerchantDetailRequest{
		MerchantID:       merchID,
		DisplayName:      "Detail Display",
		ShortDescription: "Detail Description",
	}
	created, err := s.repo.MerchantDetailCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.DisplayName, *created.DisplayName)
	detailID := int(created.MerchantDetailID)

	// 2. FindByID
	found, err := s.repo.MerchantDetailQuery.FindByID(ctx, detailID)
	s.NoError(err)
	s.Equal(created.MerchantDetailID, found.MerchantDetailID)

	// 3. FindAll
	all, err := s.repo.MerchantDetailQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 4. FindActive (before trash)
	active, err := s.repo.MerchantDetailQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 5. Update
	updateReq := &requests.UpdateMerchantDetailRequest{
		MerchantDetailID: &detailID,
		DisplayName:      "Updated Display",
		ShortDescription: "Updated Description",
	}
	updated, err := s.repo.MerchantDetailCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated Display", *updated.DisplayName)

	// 6. Trash
	trashed, err := s.repo.MerchantDetailCommand.Trash(ctx, detailID)
	s.NoError(err)
	s.NotNil(trashed)

	// 7. FindTrashed
	trashedItems, err := s.repo.MerchantDetailQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 8. FindActive after trash — should NOT include
	activeAfterTrash, err := s.repo.MerchantDetailQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(detailID, int(item.MerchantDetailID))
	}

	// 9. Restore
	restored, err := s.repo.MerchantDetailCommand.Restore(ctx, detailID)
	s.NoError(err)
	s.NotNil(restored)

	// 10. Trash again then DeletePermanent
	_, err = s.repo.MerchantDetailCommand.Trash(ctx, detailID)
	s.Require().NoError(err)

	success, err := s.repo.MerchantDetailCommand.DeletePermanent(ctx, detailID)
	s.NoError(err)
	s.True(success)

	// 11. RestoreAll
	s.repo.MerchantDetailCommand.Create(ctx, &requests.CreateMerchantDetailRequest{
		MerchantID: merchID, DisplayName: "D1", ShortDescription: "Desc1",
	})
	second, _ := s.repo.MerchantDetailCommand.Create(ctx, &requests.CreateMerchantDetailRequest{
		MerchantID: merchID, DisplayName: "D2", ShortDescription: "Desc2",
	})
	s.repo.MerchantDetailCommand.Trash(ctx, int(second.MerchantDetailID))

	resRestore, err := s.repo.MerchantDetailCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 12. DeleteAll
	s.repo.MerchantDetailCommand.Trash(ctx, int(second.MerchantDetailID))
	resDelete, err := s.repo.MerchantDetailCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestMerchantDetailRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDetailRepositoryTestSuite))
}
