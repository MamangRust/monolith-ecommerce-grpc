package product_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedCountStock = errors.ErrInternal.WithMessage("Failed to count stock")

	ErrFailedDeletingNotFoundProduct = errors.ErrNotFound.WithMessage("Product not found")
	ErrFailedDeleteImageProduct      = errors.ErrInternal.WithMessage("Failed to delete image product")

	ErrFailedFindAllProducts        = errors.ErrInternal.WithMessage("Failed to find all products")
	ErrFailedFindProductsByMerchant = errors.ErrInternal.WithMessage("Failed to find products by merchant")
	ErrFailedFindProductsByCategory = errors.ErrInternal.WithMessage("Failed to find products by category")
	ErrFailedFindProductById        = errors.ErrInternal.WithMessage("Failed to find product by ID")
	ErrFailedFindProductTrashedById = errors.ErrInternal.WithMessage("Failed to find trashed product by ID")
	ErrFailedFindProductByTrashed   = errors.ErrInternal.WithMessage("Failed to find product by trashed")

	ErrFailedFindProductsByActive  = errors.ErrInternal.WithMessage("Failed to find active products")
	ErrFailedFindProductsByTrashed = errors.ErrInternal.WithMessage("Failed to find trashed products")
	ErrFailedCreateProduct         = errors.ErrInternal.WithMessage("Failed to create product")
	ErrFailedUpdateProduct         = errors.ErrInternal.WithMessage("Failed to update product")

	ErrFailedTrashProduct               = errors.ErrInternal.WithMessage("Failed to trash product")
	ErrFailedRestoreProduct             = errors.ErrInternal.WithMessage("Failed to restore product")
	ErrFailedDeleteProductPermanent     = errors.ErrInternal.WithMessage("Failed to permanently delete product")
	ErrFailedRestoreAllProducts         = errors.ErrInternal.WithMessage("Failed to restore all products")
	ErrFailedDeleteAllProductsPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all products")
)
