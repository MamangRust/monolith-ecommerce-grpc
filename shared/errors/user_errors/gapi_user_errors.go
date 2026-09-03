package user_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcUserNotFound  = errors.NewGrpcError("User not found", http.StatusNotFound)
	ErrGrpcUserInvalidId = errors.NewGrpcError("Invalid User ID", http.StatusNotFound)

	ErrGrpcValidateCreateUser = errors.NewGrpcError("validation failed: invalid create User request", http.StatusBadRequest)
	ErrGrpcValidateUpdateUser = errors.NewGrpcError("validation failed: invalid update User request", http.StatusBadRequest)

	ErrGrpcUserInvalidEmail            = errors.NewGrpcError("Invalid email address", http.StatusBadRequest)
	ErrGrpcUserInvalidVerificationCode = errors.NewGrpcError("Invalid verification code", http.StatusBadRequest)
)
