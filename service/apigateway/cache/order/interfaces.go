package order_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderStatsCache interface {
	GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) (*response.ApiResponseOrderMonthlyTotalRevenue, bool)
	SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, data *response.ApiResponseOrderMonthlyTotalRevenue)

	GetYearlyTotalRevenueCache(ctx context.Context, year int) (*response.ApiResponseOrderYearlyTotalRevenue, bool)
	SetYearlyTotalRevenueCache(ctx context.Context, year int, data *response.ApiResponseOrderYearlyTotalRevenue)

	GetMonthlyOrderCache(ctx context.Context, year int) (*response.ApiResponseOrderMonthly, bool)
	SetMonthlyOrderCache(ctx context.Context, year int, data *response.ApiResponseOrderMonthly)

	GetYearlyOrderCache(ctx context.Context, year int) (*response.ApiResponseOrderYearly, bool)
	SetYearlyOrderCache(ctx context.Context, year int, data *response.ApiResponseOrderYearly)
}

type OrderStatsByMerchantCache interface {
	GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) (*response.ApiResponseOrderMonthlyTotalRevenue, bool)
	SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, data *response.ApiResponseOrderMonthlyTotalRevenue)

	GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) (*response.ApiResponseOrderYearlyTotalRevenue, bool)
	SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, data *response.ApiResponseOrderYearlyTotalRevenue)

	GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) (*response.ApiResponseOrderMonthly, bool)
	SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, data *response.ApiResponseOrderMonthly)

	GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) (*response.ApiResponseOrderYearly, bool)
	SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, data *response.ApiResponseOrderYearly)
}

type OrderQueryCache interface {
	GetOrderAllCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrder, bool)
	SetOrderAllCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrder)

	GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrderDeleteAt, bool)
	SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrderDeleteAt)

	GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrderDeleteAt, bool)
	SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrderDeleteAt)

	GetCachedOrderCache(ctx context.Context, order_id int) (*response.ApiResponseOrder, bool)
	SetCachedOrderCache(ctx context.Context, data *response.ApiResponseOrder)
}

type OrderCommandCache interface {
	DeleteOrderCache(ctx context.Context, orderID int)
}
