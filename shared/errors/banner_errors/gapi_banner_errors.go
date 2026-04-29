package banner_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcBannerNotFound  = errors.NewGrpcError("Banner not found", int(codes.NotFound))
	ErrGrpcBannerInvalidId = errors.NewGrpcError("Invalid Banner ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateBanner = errors.NewGrpcError("Validation failed: invalid create banner request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateBanner = errors.NewGrpcError("Validation failed: invalid update banner request", int(codes.InvalidArgument))
)
