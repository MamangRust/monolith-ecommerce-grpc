package cart_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrCartNotFoundRes   = errors.ErrNotFound.WithMessage("Cart not found")
	ErrCartAlreadyExists = errors.ErrBadRequest.WithMessage("Cart already exists")
	ErrCartInvalidData   = errors.ErrBadRequest.WithMessage("Invalid cart data")

	ErrFailedFindAllCarts = errors.ErrInternal.WithMessage("Failed to fetch carts")
	ErrFailedCreateCart   = errors.ErrInternal.WithMessage("Failed to create cart")

	ErrFailedDeleteCart     = errors.ErrInternal.WithMessage("Failed to permanently delete cart")
	ErrFailedDeleteAllCarts = errors.ErrInternal.WithMessage("Failed to permanently delete all carts")
)
