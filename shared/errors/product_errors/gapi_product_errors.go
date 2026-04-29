package product_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = errors.NewGrpcError("invalid ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateProduct = errors.NewGrpcError("validation failed: invalid create product request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateProduct = errors.NewGrpcError("validation failed: invalid update product request", int(codes.InvalidArgument))
)
