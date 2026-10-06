package orderapimapper

import (
	pborder "github.com/MamangRust/monolith-ecommerce-pb/order"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderBaseResponseMapper interface {
	ToResponseOrder(order *pborder.OrderResponse) *response.OrderResponse
	ToResponsesOrder(orders []*pborder.OrderResponse) []*response.OrderResponse
	ToApiResponseOrder(pbResponse *pborder.ApiResponseOrder) *response.ApiResponseOrder
}

type OrderQueryResponseMapper interface {
	OrderBaseResponseMapper
	ToApiResponsesOrder(pbResponse *pborder.ApiResponsesOrder) *response.ApiResponsesOrder
	ToApiResponsePaginationOrder(pbResponse *pborder.ApiResponsePaginationOrder) *response.ApiResponsePaginationOrder
	ToApiResponsePaginationOrderDeleteAt(pbResponse *pborder.ApiResponsePaginationOrderDeleteAt) *response.ApiResponsePaginationOrderDeleteAt
}

type OrderCommandResponseMapper interface {
	OrderBaseResponseMapper
	ToResponseOrderDeleteAt(order *pborder.OrderResponseDeleteAt) *response.OrderResponseDeleteAt
	ToResponsesOrderDeleteAt(orders []*pborder.OrderResponseDeleteAt) []*response.OrderResponseDeleteAt
	ToApiResponseOrderDeleteAt(pbResponse *pborder.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt
	ToApiResponseOrderDelete(pbResponse *pborder.ApiResponseOrderDelete) *response.ApiResponseOrderDelete
	ToApiResponseOrderAll(pbResponse *pborder.ApiResponseOrderAll) *response.ApiResponseOrderAll
}

type OrderStatsResponseMapper interface {
	ToOrderMonthlyPrice(category *pborder.OrderMonthlyResponse) *response.OrderMonthlyResponse
	ToOrderMonthlyPrices(c []*pborder.OrderMonthlyResponse) []*response.OrderMonthlyResponse
	ToOrderYearlyPrice(category *pborder.OrderYearlyResponse) *response.OrderYearlyResponse
	ToOrderYearlyPrices(c []*pborder.OrderYearlyResponse) []*response.OrderYearlyResponse
	ToResponseOrderMonthlyTotalRevenue(c *pborder.OrderMonthlyTotalRevenueResponse) *response.OrderMonthlyTotalRevenueResponse
	ToResponseOrderMonthlyTotalRevenues(c []*pborder.OrderMonthlyTotalRevenueResponse) []*response.OrderMonthlyTotalRevenueResponse
	ToResponseOrderYearlyTotalRevenue(c *pborder.OrderYearlyTotalRevenueResponse) *response.OrderYearlyTotalRevenueResponse
	ToResponseOrderYearlyTotalRevenues(c []*pborder.OrderYearlyTotalRevenueResponse) []*response.OrderYearlyTotalRevenueResponse

	ToApiResponseMonthlyOrder(pbResponse *pborder.ApiResponseOrderMonthly) *response.ApiResponseOrderMonthly
	ToApiResponseYearlyOrder(pbResponse *pborder.ApiResponseOrderYearly) *response.ApiResponseOrderYearly
	ToApiResponseMonthlyTotalRevenue(pbResponse *pborder.ApiResponseOrderMonthlyTotalRevenue) *response.ApiResponseOrderMonthlyTotalRevenue
	ToApiResponseYearlyTotalRevenue(pbResponse *pborder.ApiResponseOrderYearlyTotalRevenue) *response.ApiResponseOrderYearlyTotalRevenue
}
