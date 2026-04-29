package role_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrRoleNotFoundRes   = errors.ErrNotFound.WithMessage("Role not found")
	ErrFailedFindAll     = errors.ErrInternal.WithMessage("Failed to fetch Roles")
	ErrFailedFindActive  = errors.ErrInternal.WithMessage("Failed to fetch active Roles")
	ErrFailedFindTrashed = errors.ErrInternal.WithMessage("Failed to fetch trashed Roles")

	ErrFailedCreateRole = errors.ErrInternal.WithMessage("Failed to create Role")
	ErrFailedUpdateRole = errors.ErrInternal.WithMessage("Failed to update Role")

	ErrFailedTrashedRole     = errors.ErrInternal.WithMessage("Failed to move Role to trash")
	ErrFailedRestoreRole     = errors.ErrInternal.WithMessage("Failed to restore Role")
	ErrFailedDeletePermanent = errors.ErrInternal.WithMessage("Failed to delete Role permanently")

	ErrFailedRestoreAll = errors.ErrInternal.WithMessage("Failed to restore all Roles")
	ErrFailedDeleteAll  = errors.ErrInternal.WithMessage("Failed to delete all Roles permanently")
)
