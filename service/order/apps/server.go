package apps

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-order/cache"
	"github.com/MamangRust/monolith-ecommerce-order/handler"
	"github.com/MamangRust/monolith-ecommerce-order/repository"
	"github.com/MamangRust/monolith-ecommerce-order/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pborder "github.com/MamangRust/monolith-ecommerce-pb/order"
	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
	pbshipping_address "github.com/MamangRust/monolith-ecommerce-pb/shipping_address"
	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	// gRPC Client Connections
	userAddr := viper.GetString("GRPC_USER_ADDR")

	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	userQueryClient := pbuser.NewUserQueryServiceClient(userConn)

	productAddr := viper.GetString("GRPC_PRODUCT_ADDR")

	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	productQueryClient := pbproduct.NewProductQueryServiceClient(productConn)
	productCommandClient := pbproduct.NewProductCommandServiceClient(productConn)

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")
	if merchantAddr == "" {
		merchantAddr = "localhost:50055"
	}
	merchantConn, err := grpc.NewClient(merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}
	merchantQueryClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)

	orderItemAddr := viper.GetString("GRPC_ORDER_ITEM_ADDR")
	if orderItemAddr == "" {
		orderItemAddr = "localhost:50056"
	}
	orderItemConn, err := grpc.NewClient(orderItemAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order_item service: %w", err)
	}
	orderItemQueryClient := pborder_item.NewOrderItemQueryServiceClient(orderItemConn)
	orderItemCommandClient := pborder_item.NewOrderItemCommandServiceClient(orderItemConn)

	shippingAddr := viper.GetString("GRPC_SHIPPING_ADDRESS_ADDR")
	if shippingAddr == "" {
		shippingAddr = "localhost:50063"
	}
	shippingConn, err := grpc.NewClient(shippingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to shipping_address service: %w", err)
	}
	shippingCommandClient := pbshipping_address.NewShippingCommandServiceClient(shippingConn)
	shippingQueryClient := pbshipping_address.NewShippingQueryServiceClient(shippingConn)

	transactionAddr := viper.GetString("GRPC_TRANSACTION_ADDR")
	if transactionAddr == "" {
		transactionAddr = "localhost:50061"
	}
	transactionConn, err := grpc.NewClient(transactionAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to transaction service: %w", err)
	}
	transactionCommandClient := pbtransaction.NewTransactionCommandServiceClient(transactionConn)

	guardMerchant := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)
	guardProduct := resilience.NewDependencyGuard("product", 5, 30, 100, 3*time.Second, srv.Logger)
	guardOrderItem := resilience.NewDependencyGuard("order_item", 5, 30, 100, 3*time.Second, srv.Logger)
	guardUser := resilience.NewDependencyGuard("user", 5, 30, 100, 3*time.Second, srv.Logger)
	guardShipping := resilience.NewDependencyGuard("shipping_address", 5, 30, 100, 3*time.Second, srv.Logger)
	guardTransaction := resilience.NewDependencyGuard("transaction", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(&repository.Deps{
		Db:                       srv.DB,
		UserQueryClient:          userQueryClient,
		ProductQueryClient:       productQueryClient,
		ProductCommandClient:     productCommandClient,
		MerchantQueryClient:      merchantQueryClient,
		OrderItemQueryClient:     orderItemQueryClient,
		OrderItemCommandClient:   orderItemCommandClient,
		ShippingCommandClient:    shippingCommandClient,
		ShippingQueryClient:      shippingQueryClient,
		TransactionCommandClient: transactionCommandClient,
		Guards: repository.GuardOptions{
			User: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardUser),
			},
			Product: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardProduct),
			},
			Merchant: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardMerchant),
			},
			OrderItem: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardOrderItem),
			},
			Shipping: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardShipping),
			},
			Transaction: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardTransaction),
			},
		},
	})

	obs, _ := observability.NewObservability("order-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pborder.RegisterOrderQueryServiceServer(gs, h.OrderQuery)
		pborder.RegisterOrderStatsServiceServer(gs, h.OrderStats)
		pborder.RegisterOrderCommandServiceServer(gs, h.OrderCommand)
		pborder.RegisterOrderStatsByMerchantServiceServer(gs, h.OrderStatsByMerchant)
	}

	return srv, nil
}
