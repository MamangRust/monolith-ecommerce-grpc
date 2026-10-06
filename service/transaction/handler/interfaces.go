package handler

import (
	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
)

type TransactionQueryHandler interface {
	pbtransaction.TransactionQueryServiceServer
}

type TransactionCommandHandler interface {
	pbtransaction.TransactionCommandServiceServer
}

type TransactionStatsHandler interface {
	pbtransaction.TransactionStatsServiceServer
}

type TransactionStatsByMerchantHandler interface {
	pbtransaction.TransactionStatsByMerchantServiceServer
}
