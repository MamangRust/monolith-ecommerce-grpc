package shippingaddressapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type shippingAddressQueryResponseMapper struct{}

func NewShippingAddressQueryResponseMapper() ShippingAddressQueryResponseMapper {
	return &shippingAddressQueryResponseMapper{}
}

func (s *shippingAddressQueryResponseMapper) ToResponseShippingAddress(pbResponse *pb.ShippingResponse) *response.ShippingAddressResponse {
	return &response.ShippingAddressResponse{
		ID:             int(pbResponse.Id),
		OrderID:        int(pbResponse.OrderId),
		Alamat:         pbResponse.Alamat,
		Provinsi:       pbResponse.Provinsi,
		Negara:         pbResponse.Negara,
		Kota:           pbResponse.Kota,
		ShippingMethod: pbResponse.ShippingMethod,
		ShippingCost:   int(pbResponse.ShippingCost),
		CreatedAt:      pbResponse.CreatedAt,
		UpdatedAt:      pbResponse.UpdatedAt,
	}
}

func (s *shippingAddressQueryResponseMapper) ToResponsesShippingAddress(pbResponses []*pb.ShippingResponse) []*response.ShippingAddressResponse {
	var addresses []*response.ShippingAddressResponse
	for _, address := range pbResponses {
		addresses = append(addresses, s.ToResponseShippingAddress(address))
	}
	return addresses
}

func (s *shippingAddressQueryResponseMapper) ToApiResponseShippingAddress(pbResponse *pb.ApiResponseShipping) *response.ApiResponseShippingAddress {
	return &response.ApiResponseShippingAddress{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.ToResponseShippingAddress(pbResponse.Data),
	}
}

func (s *shippingAddressQueryResponseMapper) ToApiResponsesShippingAddress(pbResponse *pb.ApiResponsesShipping) *response.ApiResponsesShippingAddress {
	return &response.ApiResponsesShippingAddress{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.ToResponsesShippingAddress(pbResponse.Data),
	}
}

func (s *shippingAddressQueryResponseMapper) ToApiResponsePaginationShippingAddress(pbResponse *pb.ApiResponsePaginationShipping) *response.ApiResponsePaginationShippingAddress {
	return &response.ApiResponsePaginationShippingAddress{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.ToResponsesShippingAddress(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *shippingAddressQueryResponseMapper) ToApiResponsePaginationShippingAddressDeleteAt(pbResponse *pb.ApiResponsePaginationShippingDeleteAt) *response.ApiResponsePaginationShippingAddressDeleteAt {
	return &response.ApiResponsePaginationShippingAddressDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.ToResponsesShippingAddressDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *shippingAddressQueryResponseMapper) ToResponseShippingAddressDeleteAt(pbResponse *pb.ShippingResponseDeleteAt) *response.ShippingAddressResponseDeleteAt {
	var deletedAt *string
	if pbResponse.DeletedAt != nil {
		val := pbResponse.DeletedAt.Value
		deletedAt = &val
	}

	return &response.ShippingAddressResponseDeleteAt{
		ID:             int(pbResponse.Id),
		OrderID:        int(pbResponse.OrderId),
		Alamat:         pbResponse.Alamat,
		Provinsi:       pbResponse.Provinsi,
		Negara:         pbResponse.Negara,
		Kota:           pbResponse.Kota,
		ShippingMethod: pbResponse.ShippingMethod,
		ShippingCost:   int(pbResponse.ShippingCost),
		CreatedAt:      pbResponse.CreatedAt,
		UpdatedAt:      pbResponse.UpdatedAt,
		DeletedAt:      deletedAt,
	}
}

func (s *shippingAddressQueryResponseMapper) ToResponsesShippingAddressDeleteAt(pbResponses []*pb.ShippingResponseDeleteAt) []*response.ShippingAddressResponseDeleteAt {
	var addresses []*response.ShippingAddressResponseDeleteAt
	for _, address := range pbResponses {
		addresses = append(addresses, s.ToResponseShippingAddressDeleteAt(address))
	}
	return addresses
}
