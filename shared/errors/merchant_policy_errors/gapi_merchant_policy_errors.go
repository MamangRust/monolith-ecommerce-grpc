package merchant_policy_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcInvalidMerchantPolicyID = errors.NewGrpcError("invalid merchant policy id", http.StatusBadRequest)

	ErrGrpcValidateCreateMerchantPolicy = errors.NewGrpcError("Validation failed: invalid create merchant policy request", http.StatusBadRequest)
	ErrGrpcValidateUpdateMerchantPolicy = errors.NewGrpcError("Validation failed: invalid update merchant policy request", http.StatusBadRequest)
)
