package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/handler"
	transactionKafka "github.com/MamangRust/monolith-ecommerce-grpc-transaction/kafka"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
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

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")

	merchantConn, err := grpc.NewClient(merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}
	merchantQueryClient := pb.NewMerchantQueryServiceClient(merchantConn)

	orderAddr := viper.GetString("GRPC_ORDER_ADDR")

	orderConn, err := grpc.NewClient(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}
	orderQueryClient := pb.NewOrderQueryServiceClient(orderConn)

	orderItemAddr := viper.GetString("GRPC_ORDER_ITEM_ADDR")
	if orderItemAddr == "" {
		orderItemAddr = "order-item:50056"
	}
	orderItemConn, err := grpc.NewClient(orderItemAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order_item service: %w", err)
	}
	orderItemQueryClient := pb.NewOrderItemQueryServiceClient(orderItemConn)

	shippingAddr := viper.GetString("GRPC_SHIPPING_ADDRESS_ADDR")
	if shippingAddr == "" {
		shippingAddr = "shipping_address:50063"
	}
	shippingConn, err := grpc.NewClient(shippingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to shipping_address service: %w", err)
	}
	shippingQueryClient := pb.NewShippingQueryServiceClient(shippingConn)

	repos := repository.NewRepositories(&repository.Deps{
		DB:             srv.DB,
		UserQuery:      userQueryClient,
		MerchantQuery:  merchantQueryClient,
		OrderQuery:     orderQueryClient,
		OrderItemQuery: orderItemQueryClient,
		ShippingQuery:  shippingQueryClient,
	})
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})
	obs, _ := observability.NewObservability("transaction-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Kafka:         myKafka,
		Pool:          srv.Pool,
		Cache:         cache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	// Start the outbox relay so events committed after the transaction insert are
	// published to Kafka with durable retry and dead-letter semantics.
	go svc.Outbox.Start(srv.Ctx, service.OutboxRelayInterval, service.OutboxRelayBatchSize)

	if err := myKafka.StartConsumersWithContext(srv.Ctx, []string{"transaction-service-topic-merchant-status-event"}, "transaction-service-group", transactionKafka.NewMerchantStatusConsumer(srv.Ctx, cache.TransactionCommandCache, srv.Logger)); err != nil {
		return nil, fmt.Errorf("failed to start merchant status consumer: %w", err)
	}

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterTransactionQueryServiceServer(gs, h.TransactionQuery)
		pb.RegisterTransactionCommandServiceServer(gs, h.TransactionCommand)
		pb.RegisterTransactionStatsServiceServer(gs, h.TransactionStats)
		pb.RegisterTransactionStatsByMerchantServiceServer(gs, h.TransactionStatsByMerchant)
	}

	return srv, nil
}
