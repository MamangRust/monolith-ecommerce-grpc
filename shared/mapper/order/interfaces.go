package orderapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderBaseResponseMapper interface {
	ToResponseOrder(order *pb.OrderResponse) *response.OrderResponse
	ToResponsesOrder(orders []*pb.OrderResponse) []*response.OrderResponse
	ToApiResponseOrder(pbResponse *pb.ApiResponseOrder) *response.ApiResponseOrder
}

type OrderQueryResponseMapper interface {
	OrderBaseResponseMapper
	ToApiResponsesOrder(pbResponse *pb.ApiResponsesOrder) *response.ApiResponsesOrder
	ToApiResponsePaginationOrder(pbResponse *pb.ApiResponsePaginationOrder) *response.ApiResponsePaginationOrder
	ToApiResponsePaginationOrderDeleteAt(pbResponse *pb.ApiResponsePaginationOrderDeleteAt) *response.ApiResponsePaginationOrderDeleteAt
}

type OrderCommandResponseMapper interface {
	OrderBaseResponseMapper
	ToResponseOrderDeleteAt(order *pb.OrderResponseDeleteAt) *response.OrderResponseDeleteAt
	ToResponsesOrderDeleteAt(orders []*pb.OrderResponseDeleteAt) []*response.OrderResponseDeleteAt
	ToApiResponseOrderDeleteAt(pbResponse *pb.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt
	ToApiResponseOrderDelete(pbResponse *pb.ApiResponseOrderDelete) *response.ApiResponseOrderDelete
	ToApiResponseOrderAll(pbResponse *pb.ApiResponseOrderAll) *response.ApiResponseOrderAll
}

type OrderStatsResponseMapper interface {
	ToOrderMonthlyPrice(category *pb.OrderMonthlyResponse) *response.OrderMonthlyResponse
	ToOrderMonthlyPrices(c []*pb.OrderMonthlyResponse) []*response.OrderMonthlyResponse
	ToOrderYearlyPrice(category *pb.OrderYearlyResponse) *response.OrderYearlyResponse
	ToOrderYearlyPrices(c []*pb.OrderYearlyResponse) []*response.OrderYearlyResponse
	ToResponseOrderMonthlyTotalRevenue(c *pb.OrderMonthlyTotalRevenueResponse) *response.OrderMonthlyTotalRevenueResponse
	ToResponseOrderMonthlyTotalRevenues(c []*pb.OrderMonthlyTotalRevenueResponse) []*response.OrderMonthlyTotalRevenueResponse
	ToResponseOrderYearlyTotalRevenue(c *pb.OrderYearlyTotalRevenueResponse) *response.OrderYearlyTotalRevenueResponse
	ToResponseOrderYearlyTotalRevenues(c []*pb.OrderYearlyTotalRevenueResponse) []*response.OrderYearlyTotalRevenueResponse

	ToApiResponseMonthlyOrder(pbResponse *pb.ApiResponseOrderMonthly) *response.ApiResponseOrderMonthly
	ToApiResponseYearlyOrder(pbResponse *pb.ApiResponseOrderYearly) *response.ApiResponseOrderYearly
	ToApiResponseMonthlyTotalRevenue(pbResponse *pb.ApiResponseOrderMonthlyTotalRevenue) *response.ApiResponseOrderMonthlyTotalRevenue
	ToApiResponseYearlyTotalRevenue(pbResponse *pb.ApiResponseOrderYearlyTotalRevenue) *response.ApiResponseOrderYearlyTotalRevenue
}
