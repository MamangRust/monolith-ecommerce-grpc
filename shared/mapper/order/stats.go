package orderapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type orderStatsResponseMapper struct{}

func NewOrderStatsResponseMapper() OrderStatsResponseMapper {
	return &orderStatsResponseMapper{}
}

func (m *orderStatsResponseMapper) ToOrderMonthlyPrice(category *pb.OrderMonthlyResponse) *response.OrderMonthlyResponse {
	return &response.OrderMonthlyResponse{
		Month:          category.Month,
		OrderCount:     int(category.OrderCount),
		TotalRevenue:   int(category.TotalRevenue),
		TotalItemsSold: int(category.TotalItemsSold),
	}
}

func (m *orderStatsResponseMapper) ToOrderMonthlyPrices(c []*pb.OrderMonthlyResponse) []*response.OrderMonthlyResponse {
	var mapped []*response.OrderMonthlyResponse
	for _, item := range c {
		mapped = append(mapped, m.ToOrderMonthlyPrice(item))
	}
	return mapped
}

func (m *orderStatsResponseMapper) ToOrderYearlyPrice(category *pb.OrderYearlyResponse) *response.OrderYearlyResponse {
	return &response.OrderYearlyResponse{
		Year:               category.Year,
		OrderCount:         int(category.OrderCount),
		TotalRevenue:       int(category.TotalRevenue),
		TotalItemsSold:     int(category.TotalItemsSold),
		ActiveCashiers:     int(category.ActiveCashiers),
		UniqueProductsSold: int(category.UniqueProductsSold),
	}
}

func (m *orderStatsResponseMapper) ToOrderYearlyPrices(c []*pb.OrderYearlyResponse) []*response.OrderYearlyResponse {
	var mapped []*response.OrderYearlyResponse
	for _, item := range c {
		mapped = append(mapped, m.ToOrderYearlyPrice(item))
	}
	return mapped
}

func (m *orderStatsResponseMapper) ToResponseOrderMonthlyTotalRevenue(c *pb.OrderMonthlyTotalRevenueResponse) *response.OrderMonthlyTotalRevenueResponse {
	return &response.OrderMonthlyTotalRevenueResponse{
		Year:           c.Year,
		Month:          c.Month,
		OrderCount:     int(c.OrderCount),
		TotalRevenue:   int(c.TotalRevenue),
		TotalItemsSold: int(c.TotalItemsSold),
	}
}

func (m *orderStatsResponseMapper) ToResponseOrderMonthlyTotalRevenues(c []*pb.OrderMonthlyTotalRevenueResponse) []*response.OrderMonthlyTotalRevenueResponse {
	var mapped []*response.OrderMonthlyTotalRevenueResponse
	for _, item := range c {
		mapped = append(mapped, m.ToResponseOrderMonthlyTotalRevenue(item))
	}
	return mapped
}

func (m *orderStatsResponseMapper) ToResponseOrderYearlyTotalRevenue(c *pb.OrderYearlyTotalRevenueResponse) *response.OrderYearlyTotalRevenueResponse {
	return &response.OrderYearlyTotalRevenueResponse{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (m *orderStatsResponseMapper) ToResponseOrderYearlyTotalRevenues(c []*pb.OrderYearlyTotalRevenueResponse) []*response.OrderYearlyTotalRevenueResponse {
	var mapped []*response.OrderYearlyTotalRevenueResponse
	for _, item := range c {
		mapped = append(mapped, m.ToResponseOrderYearlyTotalRevenue(item))
	}
	return mapped
}

func (m *orderStatsResponseMapper) ToApiResponseMonthlyOrder(pbResponse *pb.ApiResponseOrderMonthly) *response.ApiResponseOrderMonthly {
	return &response.ApiResponseOrderMonthly{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToOrderMonthlyPrices(pbResponse.Data),
	}
}

func (m *orderStatsResponseMapper) ToApiResponseYearlyOrder(pbResponse *pb.ApiResponseOrderYearly) *response.ApiResponseOrderYearly {
	return &response.ApiResponseOrderYearly{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToOrderYearlyPrices(pbResponse.Data),
	}
}

func (m *orderStatsResponseMapper) ToApiResponseMonthlyTotalRevenue(pbResponse *pb.ApiResponseOrderMonthlyTotalRevenue) *response.ApiResponseOrderMonthlyTotalRevenue {
	return &response.ApiResponseOrderMonthlyTotalRevenue{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseOrderMonthlyTotalRevenues(pbResponse.Data),
	}
}

func (m *orderStatsResponseMapper) ToApiResponseYearlyTotalRevenue(pbResponse *pb.ApiResponseOrderYearlyTotalRevenue) *response.ApiResponseOrderYearlyTotalRevenue {
	return &response.ApiResponseOrderYearlyTotalRevenue{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseOrderYearlyTotalRevenues(pbResponse.Data),
	}
}
