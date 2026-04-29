package role_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcRoleNotFound  = errors.NewGrpcError("Role not found", int(codes.NotFound))
	ErrGrpcRoleInvalidId = errors.NewGrpcError("Invalid Role ID", int(codes.NotFound))

	ErrGrpcValidateCreateRole = errors.NewGrpcError("validation failed: invalid create Role request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateRole = errors.NewGrpcError("validation failed: invalid update Role request", int(codes.InvalidArgument))
)
