package transaction_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcInvalidID         = errors.NewGrpcError("invalid ID", http.StatusBadRequest)
	ErrGrpcInvalidMonth      = errors.NewGrpcError("invalid month", http.StatusBadRequest)
	ErrGrpcInvalidYear       = errors.NewGrpcError("invalid year", http.StatusBadRequest)
	ErrGrpcInvalidMerchantId = errors.NewGrpcError("invalid merchant ID", http.StatusBadRequest)

	ErrGrpcValidateCreateTransaction = errors.NewGrpcError("validation failed: invalid create transaction request", http.StatusBadRequest)
	ErrGrpcValidateUpdateTransaction = errors.NewGrpcError("validation failed: invalid update transaction request", http.StatusBadRequest)
)
