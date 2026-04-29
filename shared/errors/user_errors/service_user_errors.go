package user_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrUserIDInValid     = errors.ErrBadRequest.WithMessage("Invalid user ID")
	ErrUserNotFoundRes   = errors.ErrNotFound.WithMessage("User not found")
	ErrUserEmailAlready  = errors.ErrBadRequest.WithMessage("User email already exists")
	ErrUserPassword      = errors.ErrBadRequest.WithMessage("Failed invalid password")
	ErrFailedFindAll     = errors.ErrInternal.WithMessage("Failed to fetch users")
	ErrFailedFindActive  = errors.ErrInternal.WithMessage("Failed to fetch active users")
	ErrFailedFindTrashed = errors.ErrInternal.WithMessage("Failed to fetch trashed users")

	ErrFailedCreateUser = errors.ErrInternal.WithMessage("Failed to create user")
	ErrFailedUpdateUser = errors.ErrInternal.WithMessage("Failed to update user")

	ErrFailedTrashedUser     = errors.ErrInternal.WithMessage("Failed to move user to trash")
	ErrFailedRestoreUser     = errors.ErrInternal.WithMessage("Failed to restore user")
	ErrFailedDeletePermanent = errors.ErrInternal.WithMessage("Failed to delete user permanently")

	ErrFailedRestoreAll = errors.ErrInternal.WithMessage("Failed to restore all users")
	ErrFailedDeleteAll  = errors.ErrInternal.WithMessage("Failed to delete all users permanently")

	// ErrFailedPasswordNoMatch is returned when passwords do not match.
	ErrFailedPasswordNoMatch = errors.ErrUnauthorized.WithMessage("Failed password not match")

	// ErrAccountLocked is returned when the account is temporarily locked.
	ErrAccountLocked = errors.ErrForbidden.WithMessage("Account temporarily locked due to many failed attempts")
)
