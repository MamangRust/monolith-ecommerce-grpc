package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
}

type ProductQueryRepository interface {
	FindByID(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
}

type MerchantQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}

type ProductCommandRepository interface {
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.UpdateProductCountStockRow, error)
}

type ShippingAddressCommandRepository interface {
	Create(
		ctx context.Context,
		request *requests.CreateShippingAddressRequest,
	) (*db.CreateShippingAddressRow, error)

	Update(
		ctx context.Context,
		request *requests.UpdateShippingAddressRequest,
	) (*db.UpdateShippingAddressRow, error)

	DeleteByOrderIDPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type TransactionCommandRepository interface {
	DeleteByOrderIDPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type OrderItemQueryRepository interface {
	FindOrderItemByOrder(
		ctx context.Context,
		order_id int,
	) ([]*db.GetOrderItemsByOrderRow, error)
	CalculateTotalPrice(
		ctx context.Context,
		order_id int,
	) (*int32, error)
}

type OrderItemCommandRepository interface {
	Create(
		ctx context.Context,
		req *requests.CreateOrderItemRecordRequest,
	) (*db.CreateOrderItemRow, error)

	Update(
		ctx context.Context,
		req *requests.UpdateOrderItemRecordRequest,
	) (*db.UpdateOrderItemRow, error)

	Trash(
		ctx context.Context,
		order_id int,
	) (*db.OrderItem, error)

	Restore(
		ctx context.Context,
		order_id int,
	) (*db.OrderItem, error)

	DeletePermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)
	DeleteByOrderIDPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type OrderStatsRepository interface {
	GetMonthlyTotalRevenue(
		ctx context.Context,
		req *requests.MonthTotalRevenue,
	) ([]*db.GetMonthlyTotalRevenueRow, error)

	GetYearlyTotalRevenue(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTotalRevenueRow, error)

	GetMonthlyOrder(
		ctx context.Context,
		year int,
	) ([]*db.GetMonthlyOrderRow, error)

	GetYearlyOrder(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyOrderRow, error)
}

type OrderStatsByMerchantRepository interface {
	GetMonthlyTotalRevenueByMerchant(
		ctx context.Context,
		req *requests.MonthTotalRevenueMerchant,
	) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error)

	GetYearlyTotalRevenueByMerchant(
		ctx context.Context,
		req *requests.YearTotalRevenueMerchant,
	) ([]*db.GetYearlyTotalRevenueByMerchantRow, error)

	GetMonthlyOrderByMerchant(
		ctx context.Context,
		req *requests.MonthOrderMerchant,
	) ([]*db.GetMonthlyOrderByMerchantRow, error)

	GetYearlyOrderByMerchant(
		ctx context.Context,
		req *requests.YearOrderMerchant,
	) ([]*db.GetYearlyOrderByMerchantRow, error)
}

type OrderCommandRepository interface {
	Create(
		ctx context.Context,
		request *requests.CreateOrderRecordRequest,
	) (*db.CreateOrderRow, error)

	Update(
		ctx context.Context,
		request *requests.UpdateOrderRecordRequest,
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

type OrderQueryRepository interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersRow, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersActiveRow, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersTrashedRow, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllOrderByMerchant,
	) ([]*db.GetOrdersByMerchantRow, error)

	FindByID(
		ctx context.Context,
		order_id int,
	) (*db.GetOrderByIDRow, error)
}
