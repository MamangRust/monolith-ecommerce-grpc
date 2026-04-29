package product_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrProductNotFound     = errors.ErrNotFound.WithMessage("product not found")
	ErrFindAllProducts     = errors.ErrInternal.WithMessage("failed to find all products")
	ErrFindActiveProducts  = errors.ErrInternal.WithMessage("failed to find active products")
	ErrFindTrashedProducts = errors.ErrInternal.WithMessage("failed to find trashed products")
	ErrProductConflict     = errors.ErrConflict.WithMessage("product already exists")

	ErrCreateProduct = errors.ErrInternal.WithMessage("failed to create product")
	ErrUpdateProduct = errors.ErrInternal.WithMessage("failed to update product")

	ErrFindProductsByMerchant = errors.ErrInternal.WithMessage("failed to find products by merchant")
	ErrFindProductsByCategory = errors.ErrInternal.WithMessage("failed to find products by category")

	ErrTrashedProduct         = errors.ErrInternal.WithMessage("failed to move product to trash")
	ErrRestoreProduct         = errors.ErrInternal.WithMessage("failed to restore product from trash")
	ErrDeleteProductPermanent = errors.ErrInternal.WithMessage("failed to permanently delete product")

	ErrRestoreAllProducts = errors.ErrInternal.WithMessage("failed to restore all products")
	ErrDeleteAllProducts  = errors.ErrInternal.WithMessage("failed to permanently delete all products")

	ErrProductInternal = errors.ErrInternal.WithMessage("product internal repository error")
	ErrUpdateProductCountStock = errors.ErrInternal.WithMessage("failed to update product count stock")
)


