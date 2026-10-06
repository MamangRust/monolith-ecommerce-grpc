package transactionhandler

import (
	transaction_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/transaction"
	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/transaction"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsTransaction struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
	ApiHandler sharedErrors.ApiHandler
}

func RegisterTransactionHandler(deps *DepsTransaction) {
	mapper := apimapper.NewTransactionResponseMapper()
	statsMapper := apimapper.NewTransactionStatsResponseMapper()
	cache := transaction_cache.NewTransactionMencache(deps.CacheStore)

	queryClient := pbtransaction.NewTransactionQueryServiceClient(deps.Client)
	commandClient := pbtransaction.NewTransactionCommandServiceClient(deps.Client)
	statsClient := pbtransaction.NewTransactionStatsServiceClient(deps.Client)
	statsByMerchantClient := pbtransaction.NewTransactionStatsByMerchantServiceClient(deps.Client)

	NewTransactionQueryHandleApi(&transactionQueryHandleDeps{
		queryClient: queryClient,
		router:      deps.E,
		logger:      deps.Logger,
		mapper:      mapper.QueryMapper(),
		cache:       cache,
		apiHandler:  deps.ApiHandler,
	})

	NewTransactionCommandHandleApi(&transactionCommandHandleDeps{
		client:     commandClient,
		router:     deps.E,
		logger:     deps.Logger,
		mapper:     mapper.CommandMapper(),
		cache:      cache,
		apiHandler: deps.ApiHandler,
	})

	NewTransactionStatsHandleApi(&transactionStatsHandleDeps{
		statsClient:           statsClient,
		statsByMerchantClient: statsByMerchantClient,
		router:                deps.E,
		logger:                deps.Logger,
		statsMapper:           statsMapper,
		statsCache:            cache,
		statsByMerchantCache:  cache,
		apiHandler:            deps.ApiHandler,
	})
}
