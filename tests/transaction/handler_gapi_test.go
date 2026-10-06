package transaction_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pborder "github.com/MamangRust/monolith-ecommerce-pb/order"
	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
	pbshipping_address "github.com/MamangRust/monolith-ecommerce-pb/shipping_address"
	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	trans_cache "github.com/MamangRust/monolith-ecommerce-transaction/cache"
	trans_handler "github.com/MamangRust/monolith-ecommerce-transaction/handler"
	trans_repo "github.com/MamangRust/monolith-ecommerce-transaction/repository"
	trans_service "github.com/MamangRust/monolith-ecommerce-transaction/service"
)

type TransactionGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pbtransaction.TransactionQueryServiceClient
	commandClient pbtransaction.TransactionCommandServiceClient
}

func (s *TransactionGapiTestSuite) SetupSuite() {
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

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Transaction dependencies
	mencache := trans_cache.NewMencache(cacheStore)
	repos := trans_repo.NewRepositories(&trans_repo.Deps{
		Db:                   queries,
		UserQueryClient:      pbuser.NewUserQueryServiceClient(s.Conns["user"]),
		MerchantQueryClient:  pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		OrderQueryClient:     pborder.NewOrderQueryServiceClient(s.Conns["order"]),
		OrderItemQueryClient: pborder_item.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		ShippingQueryClient:  pbshipping_address.NewShippingQueryServiceClient(s.Conns["shipping-address"]),
	})
	svc := trans_service.NewService(&trans_service.Deps{
		Kafka:         nil,
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := trans_handler.NewHandler(&trans_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbtransaction.RegisterTransactionQueryServiceServer(server, handler.TransactionQuery)
	pbtransaction.RegisterTransactionCommandServiceServer(server, handler.TransactionCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbtransaction.NewTransactionQueryServiceClient(conn)
	s.commandClient = pbtransaction.NewTransactionCommandServiceClient(conn)
}

func (s *TransactionGapiTestSuite) TestTransactionGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, merchID, catID)
	orderID := s.SeedOrder(ctx, userID, merchID, prodID)

	// 2. Create
	createRes, err := s.commandClient.Create(ctx, &pbtransaction.CreateTransactionRequest{
		UserId:        int32(userID),
		MerchantId:    int32(merchID),
		OrderId:       int32(orderID),
		PaymentMethod: "E-Wallet",
		PaymentStatus: "pending",
		Amount:        100000,
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	transID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.queryClient.FindById(ctx, &pbtransaction.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)
	s.Equal("E-Wallet", getRes.Data.PaymentMethod)

	// 4. FindAll
	allRes, err := s.queryClient.FindAllTransactions(ctx, &pbtransaction.FindAllTransactionRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbtransaction.FindAllTransactionRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.commandClient.Update(ctx, &pbtransaction.UpdateTransactionRequest{
		TransactionId: transID,
		PaymentMethod: "Credit Card",
		PaymentStatus: "success",
		Amount:        150000,
	})
	s.Require().NoError(err)
	s.Equal("Credit Card", updateRes.Data.PaymentMethod)

	// 7. Trash
	_, err = s.commandClient.TrashedTransaction(ctx, &pbtransaction.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbtransaction.FindAllTransactionRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreTransaction(ctx, &pbtransaction.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedTransaction(ctx, &pbtransaction.FindByIdTransactionRequest{Id: transID})
	_, err = s.commandClient.DeleteTransactionPermanent(ctx, &pbtransaction.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllTransaction(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestTransactionGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionGapiTestSuite))
}
