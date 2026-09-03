package review_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcInvalidID = errors.NewGrpcError("invalid ID", http.StatusBadRequest)

	ErrGrpcValidateCreateReview = errors.NewGrpcError("validation failed: invalid create review request", http.StatusBadRequest)
	ErrGrpcValidateUpdateReview = errors.NewGrpcError("validation failed: invalid update review request", http.StatusBadRequest)
)
