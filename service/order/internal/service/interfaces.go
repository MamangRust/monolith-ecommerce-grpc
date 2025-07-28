package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderStatsService interface {
	FindMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse)
	FindYearlyTotalRevenue(ctx context.Context, year int) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse)

	FindMonthlyOrder(ctx context.Context, year int) ([]*response.OrderMonthlyResponse, *response.ErrorResponse)
	FindYearlyOrder(ctx context.Context, year int) ([]*response.OrderYearlyResponse, *response.ErrorResponse)
}

type OrderStatsByMerchantService interface {
	FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse)
	FindYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse)

	FindMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, *response.ErrorResponse)
	FindYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, *response.ErrorResponse)
}

type OrderQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, orderID int) (*response.OrderResponse, *response.ErrorResponse)
}

type OrderCommandService interface {
	CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse)
	UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse)
	TrashedOrder(ctx context.Context, orderID int) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	RestoreOrder(ctx context.Context, orderID int) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, *response.ErrorResponse)
	RestoreAllOrder(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllOrderPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
