package cart_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcCartNotFound  = errors.NewGrpcError("Cart not found", http.StatusNotFound)
	ErrGrpcCartInvalidId = errors.NewGrpcError("Invalid cart ID", http.StatusBadRequest)

	ErrGrpcFailedCreateCart      = errors.NewGrpcError("Failed to create cart", http.StatusInternalServerError)
	ErrGrpcValidateCreateCart    = errors.NewGrpcError("Validation failed: invalid create cart request", http.StatusBadRequest)
	ErrGrpcValidateDeleteCart    = errors.NewGrpcError("Validation failed: invalid delete cart request", http.StatusBadRequest)
	ErrGrpcValidateDeleteAllCart = errors.NewGrpcError("Validation failed: invalid delete all cart request", http.StatusBadRequest)
)
