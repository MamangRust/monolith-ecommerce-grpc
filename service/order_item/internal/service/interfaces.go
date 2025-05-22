package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderItemQueryService interface {
	FindAllOrderItems(req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse)
	FindOrderItemByOrder(orderID int) ([]*response.OrderItemResponse, *response.ErrorResponse)
}
