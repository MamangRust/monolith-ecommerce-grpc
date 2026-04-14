package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderStatsCache interface {
	GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) ([]*db.GetMonthlyTotalRevenueRow, bool)
	SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, data []*db.GetMonthlyTotalRevenueRow)

	GetYearlyTotalRevenueCache(ctx context.Context, year int) ([]*db.GetYearlyTotalRevenueRow, bool)
	SetYearlyTotalRevenueCache(ctx context.Context, year int, data []*db.GetYearlyTotalRevenueRow)

	GetMonthlyOrderCache(ctx context.Context, year int) ([]*db.GetMonthlyOrderRow, bool)
	SetMonthlyOrderCache(ctx context.Context, year int, data []*db.GetMonthlyOrderRow)

	GetYearlyOrderCache(ctx context.Context, year int) ([]*db.GetYearlyOrderRow, bool)
	SetYearlyOrderCache(ctx context.Context, year int, data []*db.GetYearlyOrderRow)
}

type OrderStatsByMerchantCache interface {
	GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, bool)
	SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, data []*db.GetMonthlyTotalRevenueByMerchantRow)

	GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, bool)
	SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, data []*db.GetYearlyTotalRevenueByMerchantRow)

	GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, bool)
	SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, data []*db.GetMonthlyOrderByMerchantRow)

	GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, bool)
	SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, data []*db.GetYearlyOrderByMerchantRow)
}

type OrderQueryCache interface {
	GetOrderAllCache(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersRow, *int, bool)
	SetOrderAllCache(ctx context.Context, req *requests.FindAllOrder, data []*db.GetOrdersRow, total *int)

	GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersActiveRow, *int, bool)
	SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder, data []*db.GetOrdersActiveRow, total *int)

	GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersTrashedRow, *int, bool)
	SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder, data []*db.GetOrdersTrashedRow, total *int)

	GetCachedOrderCache(ctx context.Context, orderID int) (*db.GetOrderByIDRow, bool)
	SetCachedOrderCache(ctx context.Context, data *db.GetOrderByIDRow)

	GetOrderByMerchantCache(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*db.GetOrdersByMerchantRow, *int, bool)
	SetOrderByMerchantCache(ctx context.Context, req *requests.FindAllOrderByMerchant, data []*db.GetOrdersByMerchantRow, total *int)
}

type OrderCommandCache interface {
	DeleteOrderCache(ctx context.Context, orderID int)
}
