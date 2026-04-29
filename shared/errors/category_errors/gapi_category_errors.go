package category_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateCategory = errors.NewGrpcError("Validation failed: invalid create category request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateCategory = errors.NewGrpcError("Validation failed: invalid update category request", int(codes.InvalidArgument))

	ErrGrpcCategoryNotFound          = errors.NewGrpcError("Category not found", int(codes.NotFound))
	ErrGrpcCategoryInvalidId         = errors.NewGrpcError("Invalid category ID", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidYear       = errors.NewGrpcError("Invalid year", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidMonth      = errors.NewGrpcError("Invalid month", int(codes.InvalidArgument))
	ErrGrpcCategoryInvalidMerchantId = errors.NewGrpcError("Invalid merchant ID", int(codes.InvalidArgument))

	ErrGrpcCreateCategory  = errors.NewGrpcError("Failed to create category", int(codes.Internal))
	ErrGrpcUpdateCategory  = errors.NewGrpcError("Failed to update category", int(codes.Internal))
	ErrGrpcDeleteCategory  = errors.NewGrpcError("Failed to delete category", int(codes.Internal))
	ErrGrpcFindAllCategory = errors.NewGrpcError("Failed to fetch categories", int(codes.Internal))
	ErrGrpcCategoryStats   = errors.NewGrpcError("Failed to fetch category statistics", int(codes.Internal))
)
