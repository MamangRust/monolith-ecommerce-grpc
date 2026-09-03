package merchant_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcInvalidMerchantId = errors.NewGrpcError("invalid merchant ID", http.StatusBadRequest)
	ErrGrpcMerchantInvalidID = errors.NewGrpcError("invalid ID provided", http.StatusBadRequest)

	ErrGrpcValidateCreateMerchant       = errors.NewGrpcError("Validation failed: invalid create merchant request", http.StatusBadRequest)
	ErrGrpcValidateUpdateMerchant       = errors.NewGrpcError("Validation failed: invalid update merchant request", http.StatusBadRequest)
	ErrGrpcValidateUpdateMerchantStatus = errors.NewGrpcError("Validation failed: invalid update merchant status request", http.StatusBadRequest)
	ErrGrpcFailedUpdateMerchantStatus   = errors.NewGrpcError("Failed to update merchant status", http.StatusInternalServerError)

	ErrGrpcValidateCreateMerchantDocument = errors.NewGrpcError("Validation failed: invalid create merchant document request", http.StatusBadRequest)
	ErrGrpcFailedUpdateMerchantDocument   = errors.NewGrpcError("Failed to update merchant document", http.StatusInternalServerError)
)
