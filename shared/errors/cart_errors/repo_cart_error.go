package cart_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrCartNotFound = errors.ErrNotFound.WithMessage("cart not found")
	ErrFindAllCarts = errors.ErrInternal.WithMessage("failed to find all carts")
	ErrCartConflict = errors.ErrConflict.WithMessage("cart already exists")

	ErrCreateCart = errors.ErrInternal.WithMessage("failed to create cart")

	ErrDeleteCartPermanent = errors.ErrInternal.WithMessage("failed to permanently delete cart")
	ErrDeleteAllCarts      = errors.ErrInternal.WithMessage("failed to permanently delete all carts")
)

