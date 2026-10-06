package handler

import (
	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
)

type OrderItemQueryHandler interface {
	pborder_item.OrderItemQueryServiceServer
}

type OrderItemCommandHandler interface {
	pborder_item.OrderItemCommandServiceServer
}
