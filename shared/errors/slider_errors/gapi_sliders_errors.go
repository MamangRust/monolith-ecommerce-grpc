package slider_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = errors.NewGrpcError("invalid ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateSlider = errors.NewGrpcError("validation failed: invalid create slider request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateSlider = errors.NewGrpcError("validation failed: invalid update slider request", int(codes.InvalidArgument))
)
