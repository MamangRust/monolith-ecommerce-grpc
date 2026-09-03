package merchant_business_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
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

func (s *MerchantBusinessRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *MerchantBusinessRepositoryTestSuite) TestBusinessLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	pageReq := &requests.FindAllMerchant{Search: "", Page: 1, PageSize: 10}

	// 1. Create Business Profile
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
	bizID := int(created.MerchantBusinessInfoID)

	// 2. FindByID
	found, err := s.repo.MerchantBusinessQuery.FindByID(ctx, bizID)
	s.NoError(err)
	s.Equal(created.MerchantBusinessInfoID, found.MerchantBusinessInfoID)

	// 3. FindAll
	all, err := s.repo.MerchantBusinessQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 4. FindActive (before trash)
	active, err := s.repo.MerchantBusinessQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 5. Update
	updateReq := &requests.UpdateMerchantBusinessInformationRequest{
		MerchantBusinessInfoID: &bizID,
		BusinessType:           "Wholesale",
		TaxID:                  "456-TAX",
		EstablishedYear:        2022,
		NumberOfEmployees:      25,
		WebsiteUrl:             "http://updated.mamang.corp",
	}
	updated, err := s.repo.MerchantBusinessCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Wholesale", *updated.BusinessType)

	// 6. Trash
	trashed, err := s.repo.MerchantBusinessCommand.Trash(ctx, bizID)
	s.NoError(err)
	s.NotNil(trashed)

	// 7. FindTrashed
	trashedItems, err := s.repo.MerchantBusinessQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 8. FindActive after trash — should NOT include
	activeAfterTrash, err := s.repo.MerchantBusinessQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(bizID, int(item.MerchantBusinessInfoID))
	}

	// 9. Restore
	restored, err := s.repo.MerchantBusinessCommand.Restore(ctx, bizID)
	s.NoError(err)
	s.NotNil(restored)

	// 10. Trash again then DeletePermanent
	_, err = s.repo.MerchantBusinessCommand.Trash(ctx, bizID)
	s.Require().NoError(err)

	success, err := s.repo.MerchantBusinessCommand.DeletePermanent(ctx, bizID)
	s.NoError(err)
	s.True(success)

	// 11. RestoreAll
	s.repo.MerchantBusinessCommand.Create(ctx, &requests.CreateMerchantBusinessInformationRequest{
		MerchantID: merchID, BusinessType: "T1", TaxID: "T1", EstablishedYear: 2020, NumberOfEmployees: 5,
	})
	second, _ := s.repo.MerchantBusinessCommand.Create(ctx, &requests.CreateMerchantBusinessInformationRequest{
		MerchantID: merchID, BusinessType: "T2", TaxID: "T2", EstablishedYear: 2021, NumberOfEmployees: 10,
	})
	s.repo.MerchantBusinessCommand.Trash(ctx, int(second.MerchantBusinessInfoID))

	resRestore, err := s.repo.MerchantBusinessCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 12. DeleteAll
	s.repo.MerchantBusinessCommand.Trash(ctx, int(second.MerchantBusinessInfoID))
	resDelete, err := s.repo.MerchantBusinessCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestMerchantBusinessRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantBusinessRepositoryTestSuite))
}
