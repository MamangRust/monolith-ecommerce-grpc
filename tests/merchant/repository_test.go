package merchant_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	user_repo "github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type MerchantRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo       repository.MerchantCommandRepository
	queryRepo  repository.MerchantQueryRepository
	userRepo   user_repo.UserCommandRepository
	userID     int
	merchantID int
}

func (s *MerchantRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	queries := db.New(s.DBPool())
	s.SetupUserService()
	s.repo = repository.NewMerchantCommandRepository(queries)
	s.queryRepo = repository.NewMerchantQueryRepository(queries)
	s.userRepo = user_repo.NewUserCommandRepository(queries)

	// 1. Seed dependencies
	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
}

func (s *MerchantRepositoryTestSuite) Test1_CreateMerchant() {
	req := &requests.CreateMerchantRequest{
		Name:         "Test Merchant",
		UserID:       s.userID,
		Description:  "Detailed description of the merchant.",
		Address:      "Merchant Street No. 1",
		ContactEmail: "merchant@example.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}

	merchant, err := s.repo.Create(context.Background(), req)
	s.NoError(err)
	s.NotNil(merchant)
	s.Equal(req.Name, merchant.Name)
	s.Equal(int32(s.userID), merchant.UserID)
	s.merchantID = int(merchant.MerchantID)
}

func (s *MerchantRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	found, err := s.queryRepo.FindByID(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.merchantID, int(found.MerchantID))
}

func (s *MerchantRepositoryTestSuite) Test3_UpdateMerchant() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	updateReq := &requests.UpdateMerchantRequest{
		MerchantID: &s.merchantID,
		Name:       "Updated Merchant",
		UserID:     s.userID,
		Status:     "active",
	}

	updated, err := s.repo.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(updateReq.Name, updated.Name)
	s.Equal(updateReq.Status, updated.Status)
}

func (s *MerchantRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	_, err := s.repo.Trash(ctx, s.merchantID)
	s.NoError(err)

	_, err = s.repo.Restore(ctx, s.merchantID)
	s.NoError(err)
}

func (s *MerchantRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	success, err := s.repo.DeletePermanent(ctx, s.merchantID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantRepositoryTestSuite))
}
