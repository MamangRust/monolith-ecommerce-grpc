package order_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrGetMonthlyTotalRevenue           = errors.ErrInternal.WithMessage("failed to get monthly total revenue")
	ErrGetYearlyTotalRevenue            = errors.ErrInternal.WithMessage("failed to get yearly total revenue")
	ErrGetMonthlyTotalRevenueById       = errors.ErrInternal.WithMessage("failed to get monthly total revenue by order ID")
	ErrGetYearlyTotalRevenueById        = errors.ErrInternal.WithMessage("failed to get yearly total revenue by order ID")
	ErrGetMonthlyTotalRevenueByMerchant = errors.ErrInternal.WithMessage("failed to get monthly total revenue by merchant")
	ErrGetYearlyTotalRevenueByMerchant  = errors.ErrInternal.WithMessage("failed to get yearly total revenue by merchant")

	ErrGetMonthlyOrder           = errors.ErrInternal.WithMessage("failed to get monthly orders")
	ErrGetYearlyOrder            = errors.ErrInternal.WithMessage("failed to get yearly orders")
	ErrGetMonthlyOrderByMerchant = errors.ErrInternal.WithMessage("failed to get monthly orders by merchant")
	ErrGetYearlyOrderByMerchant  = errors.ErrInternal.WithMessage("failed to get yearly orders by merchant")

	ErrFindAllOrders           = errors.ErrInternal.WithMessage("failed to find all orders")
	ErrFindByActive            = errors.ErrInternal.WithMessage("failed to find active orders")
	ErrFindByTrashed           = errors.ErrInternal.WithMessage("failed to find trashed orders")
	ErrFindByMerchant          = errors.ErrInternal.WithMessage("failed to find orders by merchant")
	ErrFindById                = errors.ErrNotFound.WithMessage("failed to find order by ID")
	ErrCreateOrder             = errors.ErrInternal.WithMessage("failed to create order")
	ErrUpdateOrder             = errors.ErrInternal.WithMessage("failed to update order")
	ErrTrashedOrder            = errors.ErrInternal.WithMessage("failed to move order to trash")
	ErrRestoreOrder            = errors.ErrInternal.WithMessage("failed to restore order from trash")
	ErrDeleteOrderPermanent    = errors.ErrInternal.WithMessage("failed to permanently delete order")
	ErrRestoreAllOrder         = errors.ErrInternal.WithMessage("failed to restore all trashed orders")
	ErrDeleteAllOrderPermanent = errors.ErrInternal.WithMessage("failed to permanently delete all trashed orders")
	ErrOrderItemNotFound       = errors.ErrNotFound.WithMessage("failed to find order items by order")

	ErrOrderNotFound = errors.ErrNotFound.WithMessage("order not found")
)


