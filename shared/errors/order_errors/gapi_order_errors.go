package order_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrGrpcInvalidYear             = errors.NewGrpcError("Invalid year", http.StatusBadRequest)
	ErrGrpcInvalidMonth            = errors.NewGrpcError("Invalid month", http.StatusBadRequest)
	ErrGrpcFailedInvalidMerchantId = errors.NewGrpcError("Invalid merchant ID", http.StatusBadRequest)
	ErrGrpcFailedInvalidId         = errors.NewGrpcError("Invalid ID", http.StatusBadRequest)

	ErrGrpcValidateCreateOrder = errors.NewGrpcError("validation failed: invalid create order request", http.StatusBadRequest)
	ErrGrpcValidateUpdateOrder = errors.NewGrpcError("validation failed: invalid update order request", http.StatusBadRequest)
)
