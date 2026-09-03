package review_detail_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcInvalidID = errors.NewGrpcError("invalid ID", http.StatusBadRequest)

	ErrGrpcValidateCreateReviewDetail = errors.NewGrpcError("Validation failed: invalid create review detail request", http.StatusBadRequest)
	ErrGrpcValidateUpdateReviewDetail = errors.NewGrpcError("Validation failed: invalid update review detail request", http.StatusBadRequest)
)
