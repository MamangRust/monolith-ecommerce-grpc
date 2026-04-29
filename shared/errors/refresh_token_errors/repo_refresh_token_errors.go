package refreshtoken_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrTokenNotFound      = errors.ErrNotFound.WithMessage("Refresh token not found")
	ErrFindByToken        = errors.ErrInternal.WithMessage("Failed to find refresh token by token")
	ErrFindByUserID       = errors.ErrInternal.WithMessage("Failed to find refresh token by user ID")
	ErrCreateRefreshToken = errors.ErrInternal.WithMessage("Failed to create refresh token")
	ErrUpdateRefreshToken = errors.ErrInternal.WithMessage("Failed to update refresh token")
	ErrDeleteRefreshToken = errors.ErrInternal.WithMessage("Failed to delete refresh token")
	ErrDeleteByUserID     = errors.ErrInternal.WithMessage("Failed to delete refresh token by user ID")
	ErrParseDate          = errors.ErrBadRequest.WithMessage("Failed to parse expiration date")
)

