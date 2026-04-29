package merchant_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantId = errors.NewGrpcError("invalid merchant ID", int(codes.InvalidArgument))
	ErrGrpcMerchantInvalidID = errors.NewGrpcError("invalid ID provided", int(codes.InvalidArgument))

	ErrGrpcValidateCreateMerchant = errors.NewGrpcError("Validation failed: invalid create merchant request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchant = errors.NewGrpcError("Validation failed: invalid update merchant request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantStatus = errors.NewGrpcError("Validation failed: invalid update merchant status request", int(codes.InvalidArgument))
	ErrGrpcFailedUpdateMerchantStatus   = errors.NewGrpcError("Failed to update merchant status", int(codes.Internal))

	ErrGrpcValidateCreateMerchantDocument = errors.NewGrpcError("Validation failed: invalid create merchant document request", int(codes.InvalidArgument))
	ErrGrpcFailedUpdateMerchantDocument   = errors.NewGrpcError("Failed to update merchant document", int(codes.Internal))
)
