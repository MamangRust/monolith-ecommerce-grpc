package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	userQueryClient := pb.NewUserQueryServiceClient(userConn)

	productAddr := viper.GetString("GRPC_PRODUCT_ADDR")
	
	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	productQueryClient := pb.NewProductQueryServiceClient(productConn)
	productCommandClient := pb.NewProductCommandServiceClient(productConn)

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")
	if merchantAddr == "" {
		merchantAddr = "50055"
	}
	merchantConn, err := grpc.NewClient("localhost:"+merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}
	merchantQueryClient := pb.NewMerchantQueryServiceClient(merchantConn)

	orderItemAddr := viper.GetString("GRPC_ORDER_ITEM_ADDR")
	if orderItemAddr == "" {
		orderItemAddr = "50056"
	}
	orderItemConn, err := grpc.NewClient("localhost:"+orderItemAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order_item service: %w", err)
	}
	orderItemQueryClient := pb.NewOrderItemQueryServiceClient(orderItemConn)
	orderItemCommandClient := pb.NewOrderItemCommandServiceClient(orderItemConn)

	shippingAddr := viper.GetString("GRPC_SHIPPING_ADDRESS_ADDR")
	if shippingAddr == "" {
		shippingAddr = "50063"
	}
	shippingConn, err := grpc.NewClient("localhost:"+shippingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to shipping_address service: %w", err)
	}
	shippingCommandClient := pb.NewShippingCommandServiceClient(shippingConn)

	transactionAddr := viper.GetString("GRPC_TRANSACTION_ADDR")
	if transactionAddr == "" {
		transactionAddr = "50061"
	}
	transactionConn, err := grpc.NewClient("localhost:"+transactionAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to transaction service: %w", err)
	}
	transactionCommandClient := pb.NewTransactionCommandServiceClient(transactionConn)

	repos := repository.NewRepositories(&repository.Deps{
		DB:               srv.DB,
		UserQuery:        userQueryClient,
		ProductQuery:     productQueryClient,
		ProductCommand:   productCommandClient,
		MerchantQuery:    merchantQueryClient,
		OrderItemQuery:   orderItemQueryClient,
		OrderItemCommand: orderItemCommandClient,
		ShippingCommand:  shippingCommandClient,
		TransactionCommand: transactionCommandClient,
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
		pb.RegisterOrderQueryServiceServer(gs, h.OrderQuery)
		pb.RegisterOrderStatsServiceServer(gs, h.OrderStats)
		pb.RegisterOrderCommandServiceServer(gs, h.OrderCommand)
	}

	return srv, nil
}
