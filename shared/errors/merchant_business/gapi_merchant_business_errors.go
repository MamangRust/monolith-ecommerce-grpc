package merchantbusiness_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateMerchantBusiness = errors.NewGrpcError("Validation failed: invalid create merchant business request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantBusiness = errors.NewGrpcError("Validation failed: invalid update merchant business request", int(codes.InvalidArgument))

	ErrGrpcMerchantBusinessNotFound  = errors.NewGrpcError("Merchant business not found", int(codes.NotFound))
	ErrGrpcInvalidMerchantBusinessId = errors.NewGrpcError("Invalid merchant business ID", int(codes.InvalidArgument))
)
