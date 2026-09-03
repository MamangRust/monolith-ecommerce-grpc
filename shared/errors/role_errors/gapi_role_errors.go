package role_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcRoleNotFound  = errors.NewGrpcError("Role not found", http.StatusNotFound)
	ErrGrpcRoleInvalidId = errors.NewGrpcError("Invalid Role ID", http.StatusNotFound)

	ErrGrpcValidateCreateRole = errors.NewGrpcError("validation failed: invalid create Role request", http.StatusBadRequest)
	ErrGrpcValidateUpdateRole = errors.NewGrpcError("validation failed: invalid update Role request", http.StatusBadRequest)
)
