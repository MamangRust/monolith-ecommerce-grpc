package auth_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var ErrGrpcLogin = errors.NewGrpcError(
	"login failed: invalid argument provided",
	int(codes.InvalidArgument),
)

var ErrGrpcGetMe = errors.NewGrpcError(
	"get user info failed: unauthenticated",
	int(codes.Unauthenticated),
)

var ErrGrpcRegisterToken = errors.NewGrpcError(
	"register failed: invalid argument",
	int(codes.InvalidArgument),
)
