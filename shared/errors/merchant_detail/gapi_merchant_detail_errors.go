package merchantdetail_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcInvalidMerchantDetailId = errors.NewGrpcError("invalid merchant detail ID", http.StatusBadRequest)

	ErrGrpcValidateCreateMerchantDetail = errors.NewGrpcError("Validation failed: invalid create merchant detail request", http.StatusBadRequest)
	ErrGrpcValidateUpdateMerchantDetail = errors.NewGrpcError("Validation failed: invalid update merchant detail request", http.StatusBadRequest)
)
