package merchantbusiness_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcValidateCreateMerchantBusiness = errors.NewGrpcError("Validation failed: invalid create merchant business request", http.StatusBadRequest)
	ErrGrpcValidateUpdateMerchantBusiness = errors.NewGrpcError("Validation failed: invalid update merchant business request", http.StatusBadRequest)

	ErrGrpcMerchantBusinessNotFound  = errors.NewGrpcError("Merchant business not found", http.StatusNotFound)
	ErrGrpcInvalidMerchantBusinessId = errors.NewGrpcError("Invalid merchant business ID", http.StatusBadRequest)
)
