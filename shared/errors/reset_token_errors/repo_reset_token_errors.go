package resettoken_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrTokenNotFound     = errors.ErrNotFound.WithMessage("Reset token not found")
	ErrCreateResetToken  = errors.ErrInternal.WithMessage("Failed to create reset token")
	ErrDeleteResetToken  = errors.ErrInternal.WithMessage("Failed to delete reset token")
	ErrDeleteByUserID    = errors.ErrInternal.WithMessage("Failed to delete reset token by user ID")
)
