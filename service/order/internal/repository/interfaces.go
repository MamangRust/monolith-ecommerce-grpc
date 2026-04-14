package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
}

type ProductQueryRepository interface {
	FindById(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}

type ProductCommandRepository interface {
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.UpdateProductCountStockRow, error)
}

type ShippingAddressCommandRepository interface {
	CreateShippingAddress(
		ctx context.Context,
		request *requests.CreateShippingAddressRequest,
	) (*db.CreateShippingAddressRow, error)

	UpdateShippingAddress(
		ctx context.Context,
		request *requests.UpdateShippingAddressRequest,
	) (*db.UpdateShippingAddressRow, error)
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
	CreateOrderItem(
		ctx context.Context,
		req *requests.CreateOrderItemRecordRequest,
	) (*db.CreateOrderItemRow, error)

	UpdateOrderItem(
		ctx context.Context,
		req *requests.UpdateOrderItemRecordRequest,
	) (*db.UpdateOrderItemRow, error)

	TrashedOrderItem(
		ctx context.Context,
		order_id int,
	) (*db.OrderItem, error)

	RestoreOrderItem(
		ctx context.Context,
		order_id int,
	) (*db.OrderItem, error)

	DeleteOrderItemPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)
	RestoreAllOrderItem(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
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
	CreateOrder(
		ctx context.Context,
		request *requests.CreateOrderRecordRequest,
	) (*db.CreateOrderRow, error)

	UpdateOrder(
		ctx context.Context,
		request *requests.UpdateOrderRecordRequest,
	) (*db.UpdateOrderRow, error)

	TrashedOrder(
		ctx context.Context,
		order_id int,
	) (*db.Order, error)

	RestoreOrder(
		ctx context.Context,
		order_id int,
	) (*db.Order, error)

	DeleteOrderPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)

	RestoreAllOrder(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}

type OrderQueryRepository interface {
	FindAllOrders(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*db.GetOrdersTrashedRow, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllOrderByMerchant,
	) ([]*db.GetOrdersByMerchantRow, error)

	FindById(
		ctx context.Context,
		order_id int,
	) (*db.GetOrderByIDRow, error)
}
