package userrole_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrAssignRoleToUser = errors.ErrInternal.WithMessage("Failed to assign role to user")
	ErrRemoveRole       = errors.ErrInternal.WithMessage("Failed to remove role from user")
)

