package transaction_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type TransactionRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *TransactionRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupShippingAddressService()
	s.SetupOrderItemService()
	s.SetupOrderService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(&repository.Deps{
		DB:             queries,
		UserQuery:      pb.NewUserQueryServiceClient(s.Conns["user"]),
		MerchantQuery:  pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		OrderQuery:     pb.NewOrderQueryServiceClient(s.Conns["order"]),
		OrderItemQuery: pb.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		ShippingQuery:  pb.NewShippingQueryServiceClient(s.Conns["shipping-address"]),
	})
}

func (s *TransactionRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *TransactionRepositoryTestSuite) TestTransactionLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, merchID, catID)
	orderID := s.SeedOrder(ctx, userID, merchID, prodID)

	pageReq := &requests.FindAllTransaction{Search: "", Page: 1, PageSize: 10}

	// 1. Create Transaction
	status := "paid"
	req := &requests.CreateTransactionRequest{
		UserID:        userID,
		OrderID:       orderID,
		MerchantID:    merchID,
		PaymentMethod: "Credit Card",
		PaymentStatus: &status,
		Amount:        100000,
	}
	created, err := s.repo.TransactionCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(orderID), created.OrderID)
	txnID := int(created.TransactionID)

	// 2. FindByID
	found, err := s.repo.TransactionQuery.FindByID(ctx, txnID)
	s.NoError(err)
	s.Equal(created.Amount, found.Amount)

	// 3. FindByOrderID
	foundByOrder, err := s.repo.TransactionQuery.FindByOrderID(ctx, orderID)
	s.NoError(err)
	s.NotNil(foundByOrder)

	// 4. FindAll
	all, err := s.repo.TransactionQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 5. FindActive (before trash)
	active, err := s.repo.TransactionQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 6. Update
	updatedStatus := "settled"
	updateReq := &requests.UpdateTransactionRequest{
		TransactionID: &txnID,
		OrderID:       orderID,
		MerchantID:    merchID,
		PaymentMethod: "Debit Card",
		Amount:        120000,
		PaymentStatus: &updatedStatus,
	}
	updated, err := s.repo.TransactionCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(120000), updated.Amount)

	// 7. Trash
	trashed, err := s.repo.TransactionCommand.Trash(ctx, txnID)
	s.NoError(err)
	s.NotNil(trashed)

	// 8. FindTrashed
	trashedItems, err := s.repo.TransactionQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 9. FindActive after trash — should NOT include
	activeAfterTrash, err := s.repo.TransactionQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(txnID, int(item.TransactionID))
	}

	// 10. Restore
	restored, err := s.repo.TransactionCommand.Restore(ctx, txnID)
	s.NoError(err)
	s.NotNil(restored)

	// 11. Trash again then DeletePermanent
	_, err = s.repo.TransactionCommand.Trash(ctx, txnID)
	s.Require().NoError(err)

	success, err := s.repo.TransactionCommand.DeletePermanent(ctx, txnID)
	s.NoError(err)
	s.True(success)

	// 12. RestoreAll
	status2 := "pending"
	second, _ := s.repo.TransactionCommand.Create(ctx, &requests.CreateTransactionRequest{
		UserID: userID, OrderID: orderID, MerchantID: merchID,
		PaymentMethod: "Bank Transfer", PaymentStatus: &status2, Amount: 50000,
	})
	s.repo.TransactionCommand.Trash(ctx, int(second.TransactionID))

	resRestore, err := s.repo.TransactionCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 13. DeleteAll
	s.repo.TransactionCommand.Trash(ctx, int(second.TransactionID))
	resDelete, err := s.repo.TransactionCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestTransactionRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
