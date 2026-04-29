package user_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrUserNotFound     = errors.ErrNotFound.WithMessage("User not found")
	ErrFindAllUsers     = errors.ErrInternal.WithMessage("Failed to find all users")
	ErrFindActiveUsers  = errors.ErrInternal.WithMessage("Failed to find active users")
	ErrFindTrashedUsers = errors.ErrInternal.WithMessage("Failed to find trashed users")
	ErrUserConflict     = errors.ErrConflict.WithMessage("User already exists")

	ErrCreateUser                 = errors.ErrInternal.WithMessage("Failed to create user")
	ErrUpdateUser                 = errors.ErrInternal.WithMessage("Failed to update user")
	ErrUpdateUserVerificationCode = errors.ErrInternal.WithMessage("Failed to update user verification code")
	ErrUpdateUserPassword         = errors.ErrInternal.WithMessage("Failed to update user password")

	ErrTrashedUser         = errors.ErrInternal.WithMessage("Failed to move user to trash")
	ErrRestoreUser         = errors.ErrInternal.WithMessage("Failed to restore user from trash")
	ErrDeleteUserPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete user")

	ErrRestoreAllUsers = errors.ErrInternal.WithMessage("Failed to restore all users")
	ErrDeleteAllUsers  = errors.ErrInternal.WithMessage("Failed to permanently delete all users")

	ErrUserInternal = errors.ErrInternal.WithMessage("User internal error")
)


