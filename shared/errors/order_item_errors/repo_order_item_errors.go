package orderitem_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrFindAllOrderItems        = errors.ErrInternal.WithMessage("failed to find all order items")
	ErrFindByActive             = errors.ErrInternal.WithMessage("failed to find active order items")
	ErrFindByTrashed            = errors.ErrInternal.WithMessage("failed to find trashed order items")
	ErrFindOrderItemByOrder     = errors.ErrNotFound.WithMessage("failed to find order items by order ID")
	ErrCalculateTotalPrice      = errors.ErrInternal.WithMessage("failed to calculate total price")
	ErrCreateOrderItem          = errors.ErrInternal.WithMessage("failed to create order item")
	ErrUpdateOrderItem          = errors.ErrInternal.WithMessage("failed to update order item")
	ErrTrashedOrderItem         = errors.ErrInternal.WithMessage("failed to move order item to trash")
	ErrRestoreOrderItem         = errors.ErrInternal.WithMessage("failed to restore order item from trash")
	ErrDeleteOrderItemPermanent = errors.ErrInternal.WithMessage("failed to permanently delete order item")
	ErrRestoreAllOrderItem      = errors.ErrInternal.WithMessage("failed to restore all trashed order items")
	ErrDeleteAllOrderPermanent  = errors.ErrInternal.WithMessage("failed to permanently delete all trashed order items")

	ErrOrderItemNotFound = errors.ErrNotFound.WithMessage("order item not found")
)

