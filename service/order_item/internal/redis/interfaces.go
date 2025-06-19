package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, bool)
	SetCachedOrderItemsAll(req *requests.FindAllOrderItems, data []*response.OrderItemResponse, total *int)

	GetCachedOrderItemActive(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, bool)
	SetCachedOrderItemActive(req *requests.FindAllOrderItems, data []*response.OrderItemResponseDeleteAt, total *int)

	GetCachedOrderItemTrashed(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, bool)
	SetCachedOrderItemTrashed(req *requests.FindAllOrderItems, data []*response.OrderItemResponseDeleteAt, total *int)

	GetCachedOrderItems(order_id int) ([]*response.OrderItemResponse, bool)
	SetCachedOrderItems(data []*response.OrderItemResponse)
}
