package refreshtoken_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrRefreshTokenNotFound = errors.ErrNotFound.WithMessage("Refresh token not found")
	ErrFailedExpire         = errors.ErrInternal.WithMessage("Failed to find refresh token by token")
	ErrFailedFindByToken    = errors.ErrInternal.WithMessage("Failed to find refresh token by token")
	ErrFailedFindByUserID   = errors.ErrInternal.WithMessage("Failed to find refresh token by user ID")
	ErrFailedInValidToken   = errors.ErrInternal.WithMessage("Failed to invalid access token")
	ErrFailedInValidUserId  = errors.ErrInternal.WithMessage("Failed to invalid user id")

	ErrFailedCreateAccess  = errors.ErrInternal.WithMessage("Failed to create access token")
	ErrFailedCreateRefresh = errors.ErrInternal.WithMessage("Failed to create refresh token")

	ErrFailedCreateRefreshToken  = errors.ErrInternal.WithMessage("Failed to create refresh token")
	ErrFailedUpdateRefreshToken  = errors.ErrInternal.WithMessage("Failed to update refresh token")
	ErrFailedDeleteRefreshToken  = errors.ErrInternal.WithMessage("Failed to delete refresh token")
	ErrFailedDeleteByUserID      = errors.ErrInternal.WithMessage("Failed to delete refresh token by user ID")
	ErrFailedParseExpirationDate = errors.ErrBadRequest.WithMessage("Failed to parse expiration date")
)
