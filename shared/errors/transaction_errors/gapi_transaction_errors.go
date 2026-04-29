package transaction_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID         = errors.NewGrpcError("invalid ID", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth      = errors.NewGrpcError("invalid month", int(codes.InvalidArgument))
	ErrGrpcInvalidYear       = errors.NewGrpcError("invalid year", int(codes.InvalidArgument))
	ErrGrpcInvalidMerchantId = errors.NewGrpcError("invalid merchant ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateTransaction = errors.NewGrpcError("validation failed: invalid create transaction request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateTransaction = errors.NewGrpcError("validation failed: invalid update transaction request", int(codes.InvalidArgument))
)
