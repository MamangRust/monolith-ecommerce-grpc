package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type TransactionQueryHandler interface {
	pb.TransactionQueryServiceServer
}

type TransactionCommandHandler interface {
	pb.TransactionCommandServiceServer
}

type TransactionStatsHandler interface {
	pb.TransactionStatsServiceServer
}

type TransactionStatsByMerchantHandler interface {
	pb.TransactionStatsByMerchantServiceServer
}
