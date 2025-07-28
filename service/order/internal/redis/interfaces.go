package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderStatsCache interface {
	GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, bool)
	SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, data []*response.OrderMonthlyTotalRevenueResponse)

	GetYearlyTotalRevenueCache(ctx context.Context, year int) ([]*response.OrderYearlyTotalRevenueResponse, bool)
	SetYearlyTotalRevenueCache(ctx context.Context, year int, data []*response.OrderYearlyTotalRevenueResponse)

	GetMonthlyOrderCache(ctx context.Context, year int) ([]*response.OrderMonthlyResponse, bool)
	SetMonthlyOrderCache(ctx context.Context, year int, data []*response.OrderMonthlyResponse)

	GetYearlyOrderCache(ctx context.Context, year int) ([]*response.OrderYearlyResponse, bool)
	SetYearlyOrderCache(ctx context.Context, year int, data []*response.OrderYearlyResponse)
}

type OrderStatsByMerchantCache interface {
	GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, bool)
	SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, data []*response.OrderMonthlyTotalRevenueResponse)

	GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, bool)
	SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, data []*response.OrderYearlyTotalRevenueResponse)

	GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, bool)
	SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, data []*response.OrderMonthlyResponse)

	GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, bool)
	SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, data []*response.OrderYearlyResponse)
}

type OrderQueryCache interface {
	GetOrderAllCache(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponse, *int, bool)
	SetOrderAllCache(ctx context.Context, req *requests.FindAllOrder, data []*response.OrderResponse, total *int)

	GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, bool)
	SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder, data []*response.OrderResponseDeleteAt, total *int)

	GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, bool)
	SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder, data []*response.OrderResponseDeleteAt, total *int)

	GetCachedOrderCache(ctx context.Context, orderID int) (*response.OrderResponse, bool)
	SetCachedOrderCache(ctx context.Context, data *response.OrderResponse)
}

type OrderCommandCache interface {
	DeleteOrderCache(ctx context.Context, orderID int)
}
