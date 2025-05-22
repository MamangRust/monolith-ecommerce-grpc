package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type TransactionHandleGrpc interface {
	pb.TransactionServiceServer
}
