package merchantaward_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcMerchantInvalidId = errors.NewGrpcError("Invalid merchant ID", http.StatusBadRequest)

	ErrGrpcValidateCreateMerchantAward = errors.NewGrpcError("Validation failed: invalid create merchant award request", http.StatusBadRequest)
	ErrGrpcValidateUpdateMerchantAward = errors.NewGrpcError("Validation failed: invalid update merchant award request", http.StatusBadRequest)
)
