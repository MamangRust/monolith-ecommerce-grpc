package refreshtoken_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var ErrGrpcRefreshToken = errors.NewGrpcError("refresh token failed", int(codes.Unauthenticated))
