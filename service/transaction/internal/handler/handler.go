package handler

import "github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/service"

type Deps struct {
	Service service.Service
}

type Handler struct {
	Transaction TransactionHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Transaction: NewTransactionHandleGrpc(deps.Service),
	}
}
