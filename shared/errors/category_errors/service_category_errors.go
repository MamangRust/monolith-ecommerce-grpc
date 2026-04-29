package category_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedFindMonthlyTotalPrice           = errors.ErrInternal.WithMessage("Failed to find monthly total price")
	ErrFailedFindYearlyTotalPrice            = errors.ErrInternal.WithMessage("Failed to find yearly total price")
	ErrFailedFindMonthlyTotalPriceById       = errors.ErrInternal.WithMessage("Failed to find monthly total price by category ID")
	ErrFailedFindYearlyTotalPriceById        = errors.ErrInternal.WithMessage("Failed to find yearly total price by category ID")
	ErrFailedFindMonthlyTotalPriceByMerchant = errors.ErrInternal.WithMessage("Failed to find monthly total price by merchant")
	ErrFailedFindYearlyTotalPriceByMerchant  = errors.ErrInternal.WithMessage("Failed to find yearly total price by merchant")

	ErrFailedFindMonthPrice           = errors.ErrInternal.WithMessage("Failed to find monthly price")
	ErrFailedFindYearPrice            = errors.ErrInternal.WithMessage("Failed to find yearly price")
	ErrFailedFindMonthPriceByMerchant = errors.ErrInternal.WithMessage("Failed to find monthly price by merchant")
	ErrFailedFindYearPriceByMerchant  = errors.ErrInternal.WithMessage("Failed to find yearly price by merchant")
	ErrFailedFindMonthPriceById       = errors.ErrInternal.WithMessage("Failed to find monthly price by category ID")
	ErrFailedFindYearPriceById        = errors.ErrInternal.WithMessage("Failed to find yearly price by category ID")

	ErrFailedFindAllCategories     = errors.ErrInternal.WithMessage("Failed to find all categories")
	ErrFailedFindActiveCategories  = errors.ErrInternal.WithMessage("Failed to find active categories")
	ErrFailedFindTrashedCategories = errors.ErrInternal.WithMessage("Failed to find trashed categories")
	ErrFailedFindCategoryById      = errors.ErrInternal.WithMessage("Failed to find category by ID")
	ErrFailedFindCategoryIdTrashed = errors.ErrInternal.WithMessage("Failed to find category ID trashed")
	ErrFailedRemoveImageCategory   = errors.ErrInternal.WithMessage("Failed to remove image category")

	ErrFailedCreateCategory               = errors.ErrInternal.WithMessage("Failed to create category")
	ErrFailedUpdateCategory               = errors.ErrInternal.WithMessage("Failed to update category")
	ErrFailedTrashedCategory              = errors.ErrInternal.WithMessage("Failed to move category to trash")
	ErrFailedRestoreCategory              = errors.ErrInternal.WithMessage("Failed to restore category")
	ErrFailedDeleteCategoryPermanent      = errors.ErrInternal.WithMessage("Failed to permanently delete category")
	ErrFailedRestoreAllCategories         = errors.ErrInternal.WithMessage("Failed to restore all categories")
	ErrFailedDeleteAllCategoriesPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all categories")
)
