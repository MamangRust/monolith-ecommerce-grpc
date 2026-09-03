package category_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcValidateCreateCategory = errors.NewGrpcError("Validation failed: invalid create category request", http.StatusBadRequest)
	ErrGrpcValidateUpdateCategory = errors.NewGrpcError("Validation failed: invalid update category request", http.StatusBadRequest)

	ErrGrpcCategoryNotFound          = errors.NewGrpcError("Category not found", http.StatusNotFound)
	ErrGrpcCategoryInvalidId         = errors.NewGrpcError("Invalid category ID", http.StatusBadRequest)
	ErrGrpcCategoryInvalidYear       = errors.NewGrpcError("Invalid year", http.StatusBadRequest)
	ErrGrpcCategoryInvalidMonth      = errors.NewGrpcError("Invalid month", http.StatusBadRequest)
	ErrGrpcCategoryInvalidMerchantId = errors.NewGrpcError("Invalid merchant ID", http.StatusBadRequest)

	ErrGrpcCreateCategory  = errors.NewGrpcError("Failed to create category", http.StatusInternalServerError)
	ErrGrpcUpdateCategory  = errors.NewGrpcError("Failed to update category", http.StatusInternalServerError)
	ErrGrpcDeleteCategory  = errors.NewGrpcError("Failed to delete category", http.StatusInternalServerError)
	ErrGrpcFindAllCategory = errors.NewGrpcError("Failed to fetch categories", http.StatusInternalServerError)
	ErrGrpcCategoryStats   = errors.NewGrpcError("Failed to fetch category statistics", http.StatusInternalServerError)
)
