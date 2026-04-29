package order_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedNotDeleteAtOrder = errors.ErrInternal.WithMessage("Failed to delete at order")

	ErrInsufficientProductStock  = errors.ErrBadRequest.WithMessage("Insufficient product stock")
	ErrFailedInvalidCountInStock = errors.ErrInternal.WithMessage("Failed to find invalid count in stock")

	ErrFailedFindMonthlyTotalRevenue           = errors.ErrInternal.WithMessage("Failed to find monthly total revenue")
	ErrFailedFindYearlyTotalRevenue            = errors.ErrInternal.WithMessage("Failed to find yearly total revenue")
	ErrFailedFindMonthlyTotalRevenueById       = errors.ErrInternal.WithMessage("Failed to find monthly total revenue by order ID")
	ErrFailedFindYearlyTotalRevenueById        = errors.ErrInternal.WithMessage("Failed to find yearly total revenue by order ID")
	ErrFailedFindMonthlyTotalRevenueByMerchant = errors.ErrInternal.WithMessage("Failed to find monthly total revenue by merchant")
	ErrFailedFindYearlyTotalRevenueByMerchant  = errors.ErrInternal.WithMessage("Failed to find yearly total revenue by merchant")

	ErrFailedFindMonthlyOrder           = errors.ErrInternal.WithMessage("Failed to find monthly order")
	ErrFailedFindYearlyOrder            = errors.ErrInternal.WithMessage("Failed to find yearly order")
	ErrFailedFindMonthlyOrderByMerchant = errors.ErrInternal.WithMessage("Failed to find monthly order by merchant")
	ErrFailedFindYearlyOrderByMerchant  = errors.ErrInternal.WithMessage("Failed to find yearly order by merchant")

	ErrFailedFindAllOrders           = errors.ErrInternal.WithMessage("Failed to find all orders")
	ErrFailedFindOrderById           = errors.ErrInternal.WithMessage("Failed to find order by ID")
	ErrFailedFindOrdersByActive      = errors.ErrInternal.WithMessage("Failed to find active orders")
	ErrFailedFindOrdersByTrashed     = errors.ErrInternal.WithMessage("Failed to find trashed orders")
	ErrFailedFindOrdersByMerchant    = errors.ErrInternal.WithMessage("Failed to find orders by merchant")
	ErrFailedCreateOrder             = errors.ErrInternal.WithMessage("Failed to create order")
	ErrFailedUpdateOrder             = errors.ErrInternal.WithMessage("Failed to update order")
	ErrFailedTrashOrder              = errors.ErrInternal.WithMessage("Failed to trash order")
	ErrFailedRestoreOrder            = errors.ErrInternal.WithMessage("Failed to restore order")
	ErrFailedDeleteOrderPermanent    = errors.ErrInternal.WithMessage("Failed to permanently delete order")
	ErrFailedRestoreAllOrder         = errors.ErrInternal.WithMessage("Failed to restore all orders")
	ErrFailedDeleteAllOrderPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all orders")
)
