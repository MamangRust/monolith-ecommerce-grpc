package orderitem_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedOrderItemEmptyId     = errors.ErrBadRequest.WithMessage("Order item ID is empty")
	ErrFailedNotDeleteAtOrderItem = errors.ErrInternal.WithMessage("Failed to delete at order item")
	ErrFailedOrderItemEmpty       = errors.ErrInternal.WithMessage("Failed to find order item")
	ErrFailedInvalidQuantity      = errors.ErrBadRequest.WithMessage("Invalid quantity")

	ErrFailedOrderItemNotFound = errors.ErrNotFound.WithMessage("Order item not found")
	ErrFailedTrashedOrderItem  = errors.ErrBadRequest.WithMessage("Order item is already trashed")
	ErrFailedRestoreOrderItem  = errors.ErrInternal.WithMessage("Failed to restore order item")
	ErrFailedDeleteOrderItem   = errors.ErrInternal.WithMessage("Failed to delete order item")

	ErrFailedCreateOrderItem = errors.ErrInternal.WithMessage("Failed to create order item")
	ErrFailedUpdateOrderItem = errors.ErrInternal.WithMessage("Failed to update order item")
	ErrFailedCalculateTotal  = errors.ErrInternal.WithMessage("Failed to calculate total")

	ErrFailedFindAllOrderItems       = errors.ErrInternal.WithMessage("Failed to find all order items")
	ErrFailedFindOrderItemsByActive  = errors.ErrInternal.WithMessage("Failed to find active order items")
	ErrFailedFindOrderItemsByTrashed = errors.ErrInternal.WithMessage("Failed to find trashed order items")
	ErrFailedFindOrderItemByOrder    = errors.ErrInternal.WithMessage("Failed to find order items by order ID")

	ErrFailedRestoreAllOrderItem = errors.ErrInternal.WithMessage("Failed to restore all order items")
	ErrFailedDeleteAllOrderItem  = errors.ErrInternal.WithMessage("Failed to delete all order items")
)
