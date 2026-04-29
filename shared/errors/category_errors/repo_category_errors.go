package category_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrGetMonthlyTotalPrice           = errors.ErrInternal.WithMessage("failed to get monthly total price for categories")
	ErrGetYearlyTotalPrices           = errors.ErrInternal.WithMessage("failed to get yearly total prices for categories")
	ErrGetMonthlyTotalPriceById       = errors.ErrInternal.WithMessage("failed to get monthly total price by category ID")
	ErrGetYearlyTotalPricesById       = errors.ErrInternal.WithMessage("failed to get yearly total prices by category ID")
	ErrGetMonthlyTotalPriceByMerchant = errors.ErrInternal.WithMessage("failed to get monthly total price by merchant")
	ErrGetYearlyTotalPricesByMerchant = errors.ErrInternal.WithMessage("failed to get yearly total prices by merchant")

	ErrGetMonthPrice           = errors.ErrInternal.WithMessage("failed to get month price for categories")
	ErrGetYearPrice            = errors.ErrInternal.WithMessage("failed to get year price for categories")
	ErrGetMonthPriceByMerchant = errors.ErrInternal.WithMessage("failed to get month price by merchant")
	ErrGetYearPriceByMerchant  = errors.ErrInternal.WithMessage("failed to get year price by merchant")
	ErrGetMonthPriceById       = errors.ErrInternal.WithMessage("failed to get month price by category ID")
	ErrGetYearPriceById        = errors.ErrInternal.WithMessage("failed to get year price by category ID")

	ErrFindAllCategory         = errors.ErrInternal.WithMessage("failed to find all categories")
	ErrFindByActiveCategory    = errors.ErrInternal.WithMessage("failed to find active categories")
	ErrFindByTrashedCategory   = errors.ErrInternal.WithMessage("failed to find trashed categories")
	ErrFindCategoryById        = errors.ErrInternal.WithMessage("failed to find category by ID")
	ErrFindCategoryByIdTrashed = errors.ErrInternal.WithMessage("failed to find trashed category by ID")
	ErrCategoryNotFound        = errors.ErrNotFound.WithMessage("Category not found")

	ErrCreateCategory            = errors.ErrInternal.WithMessage("failed to create category")
	ErrUpdateCategory            = errors.ErrInternal.WithMessage("failed to update category")
	ErrTrashedCategory           = errors.ErrInternal.WithMessage("failed to move category to trash")
	ErrRestoreCategory           = errors.ErrInternal.WithMessage("failed to restore category")
	ErrDeleteCategoryPermanently = errors.ErrInternal.WithMessage("failed to permanently delete category")

	ErrRestoreAllCategories         = errors.ErrInternal.WithMessage("failed to restore all categories")
	ErrDeleteAllPermanentCategories = errors.ErrInternal.WithMessage("failed to permanently delete all categories")
)

