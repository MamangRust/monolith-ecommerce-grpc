package transactionhandler

import (
	transaction_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/cache/transaction"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/transaction"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type DepsTransaction struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterTransactionHandler(deps *DepsTransaction) {
	mapper := apimapper.NewTransactionResponseMapper()
	statsMapper := apimapper.NewTransactionStatsResponseMapper()
	cache := transaction_cache.NewTransactionMencache(deps.CacheStore)

	queryClient := pb.NewTransactionQueryServiceClient(deps.Client)
	commandClient := pb.NewTransactionCommandServiceClient(deps.Client)
	statsClient := pb.NewTransactionStatsServiceClient(deps.Client)
	statsByMerchantClient := pb.NewTransactionStatsByMerchantServiceClient(deps.Client)

	NewTransactionQueryHandleApi(&transactionQueryHandleDeps{
		queryClient: queryClient,
		router:      deps.E,
		logger:      deps.Logger,
		mapper:      mapper.QueryMapper(),
		cache:       cache,
	})

	NewTransactionCommandHandleApi(&transactionCommandHandleDeps{
		client: commandClient,
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache,
	})

	NewTransactionStatsHandleApi(&transactionStatsHandleDeps{
		statsClient:           statsClient,
		statsByMerchantClient: statsByMerchantClient,
		router:                deps.E,
		logger:                deps.Logger,
		statsMapper:           statsMapper,
		statsCache:            cache,
		statsByMerchantCache:  cache,
	})
}
