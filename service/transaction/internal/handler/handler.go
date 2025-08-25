package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Transaction TransactionHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Transaction: NewTransactionHandleGrpc(deps.Service, deps.Logger),
	}
}
