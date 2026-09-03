package banner_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcBannerNotFound  = errors.NewGrpcError("Banner not found", http.StatusNotFound)
	ErrGrpcBannerInvalidId = errors.NewGrpcError("Invalid Banner ID", http.StatusBadRequest)

	ErrGrpcValidateCreateBanner = errors.NewGrpcError("Validation failed: invalid create banner request", http.StatusBadRequest)
	ErrGrpcValidateUpdateBanner = errors.NewGrpcError("Validation failed: invalid update banner request", http.StatusBadRequest)
)
