package auth_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var ErrGrpcLogin = errors.NewGrpcError(
	"login failed: invalid argument provided",
	http.StatusBadRequest,
)

var ErrGrpcGetMe = errors.NewGrpcError(
	"get user info failed: unauthenticated",
	http.StatusUnauthorized,
)

var ErrGrpcRegisterToken = errors.NewGrpcError(
	"register failed: invalid argument",
	http.StatusBadRequest,
)
