package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderStatsService interface {
	FindMonthlyTotalRevenue(
		ctx context.Context,
		req *requests.MonthTotalRevenue,
	) ([]*db.GetMonthlyTotalRevenueRow, error)

	FindYearlyTotalRevenue(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTotalRevenueRow, error)
	FindMonthlyOrder(
		ctx context.Context,
		year int,
	) ([]*db.GetMonthlyOrderRow, error)

	FindYearlyOrder(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyOrderRow, error)
}

type OrderStatsByMerchantService interface {
	FindMonthlyTotalRevenueByMerchant(
		ctx context.Context,
		req *requests.MonthTotalRevenueMerchant,
	) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error)

	FindYearlyTotalRevenueByMerchant(
		ctx context.Context,
		req *requests.YearTotalRevenueMerchant,
	) ([]*db.GetYearlyTotalRevenueByMerchantRow, error)

	FindMonthlyOrderByMerchant(
		ctx context.Context,
		req *requests.MonthOrderMerchant,
	) ([]*db.GetMonthlyOrderByMerchantRow, error)

	FindYearlyOrderByMerchant(
		ctx context.Context,
		req *requests.YearOrderMerchant,
	) ([]*db.GetYearlyOrderByMerchantRow, error)
}

type OrderQueryService interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersRow, *int, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersActiveRow, *int, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersTrashedRow, *int, error)

	FindByID(
		ctx context.Context,
		order_id int,
	) (*db.GetOrderByIDRow, error)
}

type OrderCommandService interface {
	Create(
		ctx context.Context,
		request *requests.CreateOrderRequest,
	) (*db.CreateOrderRow, error)

	Update(
		ctx context.Context,
		request *requests.UpdateOrderRequest,
	) (*db.UpdateOrderRow, error)

	Trash(
		ctx context.Context,
		order_id int,
	) (*db.Order, error)

	Restore(
		ctx context.Context,
		order_id int,
	) (*db.Order, error)

	DeletePermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
