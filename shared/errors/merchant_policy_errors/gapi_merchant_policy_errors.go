package merchant_policy_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantPolicyID = errors.NewGrpcError("invalid merchant policy id", int(codes.InvalidArgument))

	ErrGrpcValidateCreateMerchantPolicy = errors.NewGrpcError("Validation failed: invalid create merchant policy request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantPolicy = errors.NewGrpcError("Validation failed: invalid update merchant policy request", int(codes.InvalidArgument))
)
