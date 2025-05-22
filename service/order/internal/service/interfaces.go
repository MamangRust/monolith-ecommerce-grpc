package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderStatsService interface {
	FindMonthlyTotalRevenue(req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse)
	FindYearlyTotalRevenue(year int) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse)

	FindMonthlyOrder(year int) ([]*response.OrderMonthlyResponse, *response.ErrorResponse)
	FindYearlyOrder(year int) ([]*response.OrderYearlyResponse, *response.ErrorResponse)
}

type OrderStatsByMerchantService interface {
	FindMonthlyTotalRevenueByMerchant(req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse)
	FindYearlyTotalRevenueByMerchant(req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse)

	FindMonthlyOrderByMerchant(req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, *response.ErrorResponse)
	FindYearlyOrderByMerchant(req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, *response.ErrorResponse)
}

type OrderQueryService interface {
	FindAll(req *requests.FindAllOrder) ([]*response.OrderResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(order_id int) (*response.OrderResponse, *response.ErrorResponse)
}

type OrderCommandService interface {
	CreateOrder(req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse)
	UpdateOrder(req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse)
	TrashedOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	RestoreOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	DeleteOrderPermanent(order_id int) (bool, *response.ErrorResponse)
	RestoreAllOrder() (bool, *response.ErrorResponse)
	DeleteAllOrderPermanent() (bool, *response.ErrorResponse)
}
