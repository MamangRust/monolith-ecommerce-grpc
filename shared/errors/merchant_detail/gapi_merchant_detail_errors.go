package merchantdetail_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantDetailId = errors.NewGrpcError("invalid merchant detail ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateMerchantDetail = errors.NewGrpcError("Validation failed: invalid create merchant detail request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantDetail = errors.NewGrpcError("Validation failed: invalid update merchant detail request", int(codes.InvalidArgument))
)
