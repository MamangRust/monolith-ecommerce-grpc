package role_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrRoleNotFound     = errors.ErrNotFound.WithMessage("Role not found")
	ErrFindAllRoles     = errors.ErrInternal.WithMessage("Failed to find all roles")
	ErrFindActiveRoles  = errors.ErrInternal.WithMessage("Failed to find active roles")
	ErrFindTrashedRoles = errors.ErrInternal.WithMessage("Failed to find trashed roles")
	ErrRoleConflict     = errors.ErrConflict.WithMessage("Role already exists")

	ErrCreateRole = errors.ErrInternal.WithMessage("Failed to create role")
	ErrUpdateRole = errors.ErrInternal.WithMessage("Failed to update role")

	ErrTrashedRole         = errors.ErrInternal.WithMessage("Failed to move role to trash")
	ErrRestoreRole         = errors.ErrInternal.WithMessage("Failed to restore role from trash")
	ErrDeleteRolePermanent = errors.ErrInternal.WithMessage("Failed to permanently delete role")

	ErrRestoreAllRoles = errors.ErrInternal.WithMessage("Failed to restore all roles")
	ErrDeleteAllRoles  = errors.ErrInternal.WithMessage("Failed to permanently delete all roles")

	ErrAssignRole = errors.ErrInternal.WithMessage("Failed to assign role to user")
	ErrRemoveRole = errors.ErrInternal.WithMessage("Failed to remove role from user")
)

