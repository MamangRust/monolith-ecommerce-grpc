package orderapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type orderQueryResponseMapper struct{}

func NewOrderQueryResponseMapper() OrderQueryResponseMapper {
	return &orderQueryResponseMapper{}
}

func (o *orderQueryResponseMapper) ToResponseOrder(order *pb.OrderResponse) *response.OrderResponse {
	if order == nil { return nil }
	return &response.OrderResponse{
		ID:         int(order.Id),
		MerchantID: int(order.MerchantId),
		UserID:     int(order.UserId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderQueryResponseMapper) ToResponsesOrder(orders []*pb.OrderResponse) []*response.OrderResponse {
	var mappedOrders []*response.OrderResponse
	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrder(order))
	}
	return mappedOrders
}

func (o *orderQueryResponseMapper) ToApiResponseOrder(pbResponse *pb.ApiResponseOrder) *response.ApiResponseOrder {
	return &response.ApiResponseOrder{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrder(pbResponse.Data),
	}
}

func (o *orderQueryResponseMapper) ToApiResponsesOrder(pbResponse *pb.ApiResponsesOrder) *response.ApiResponsesOrder {
	return &response.ApiResponsesOrder{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponsesOrder(pbResponse.Data),
	}
}

func (o *orderQueryResponseMapper) ToApiResponsePaginationOrder(pbResponse *pb.ApiResponsePaginationOrder) *response.ApiResponsePaginationOrder {
	return &response.ApiResponsePaginationOrder{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrder(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (o *orderQueryResponseMapper) ToApiResponsePaginationOrderDeleteAt(pbResponse *pb.ApiResponsePaginationOrderDeleteAt) *response.ApiResponsePaginationOrderDeleteAt {
	var mappedOrders []*response.OrderResponseDeleteAt
	for _, order := range pbResponse.Data {
		var deletedAt string
		if order.DeletedAt != nil {
			deletedAt = order.DeletedAt.Value
		}
		mappedOrders = append(mappedOrders, &response.OrderResponseDeleteAt{
			ID:         int(order.Id),
			MerchantID: int(order.MerchantId),
			UserID:     int(order.UserId),
			TotalPrice: int(order.TotalPrice),
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
			DeletedAt:  &deletedAt,
		})
	}

	return &response.ApiResponsePaginationOrderDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       mappedOrders,
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
