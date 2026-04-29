package cart_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcCartNotFound  = errors.NewGrpcError("Cart not found", int(codes.NotFound))
	ErrGrpcCartInvalidId = errors.NewGrpcError("Invalid cart ID", int(codes.InvalidArgument))

	ErrGrpcFailedCreateCart   = errors.NewGrpcError("Failed to create cart", int(codes.Internal))
	ErrGrpcValidateCreateCart = errors.NewGrpcError("Validation failed: invalid create cart request", int(codes.InvalidArgument))
	ErrGrpcValidateDeleteCart = errors.NewGrpcError("Validation failed: invalid delete cart request", int(codes.InvalidArgument))
	ErrGrpcValidateDeleteAllCart = errors.NewGrpcError("Validation failed: invalid delete all cart request", int(codes.InvalidArgument))
)
