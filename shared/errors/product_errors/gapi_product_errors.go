package product_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcInvalidID = errors.NewGrpcError("invalid ID", http.StatusBadRequest)

	ErrGrpcValidateCreateProduct = errors.NewGrpcError("validation failed: invalid create product request", http.StatusBadRequest)
	ErrGrpcValidateUpdateProduct = errors.NewGrpcError("validation failed: invalid update product request", http.StatusBadRequest)
)
