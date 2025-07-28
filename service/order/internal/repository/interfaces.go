package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindById(ctx context.Context, userID int) (*record.UserRecord, error)
}

type ProductQueryRepository interface {
	FindById(ctx context.Context, productID int) (*record.ProductRecord, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, merchantID int) (*record.MerchantRecord, error)
}

type ProductCommandRepository interface {
	UpdateProductCountStock(ctx context.Context, productID int, stock int) (*record.ProductRecord, error)
}

type ShippingAddressCommandRepository interface {
	CreateShippingAddress(ctx context.Context, request *requests.CreateShippingAddressRequest) (*record.ShippingAddressRecord, error)
	UpdateShippingAddress(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*record.ShippingAddressRecord, error)
}

type OrderItemQueryRepository interface {
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*record.OrderItemRecord, error)
	CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error)
}

type OrderItemCommandRepository interface {
	CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*record.OrderItemRecord, error)
	UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*record.OrderItemRecord, error)
	TrashedOrderItem(ctx context.Context, orderID int) (*record.OrderItemRecord, error)
	RestoreOrderItem(ctx context.Context, orderID int) (*record.OrderItemRecord, error)
	DeleteOrderItemPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAllOrderItem(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}

type OrderStatsRepository interface {
	GetMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*record.OrderMonthlyTotalRevenueRecord, error)
	GetYearlyTotalRevenue(ctx context.Context, year int) ([]*record.OrderYearlyTotalRevenueRecord, error)

	GetMonthlyOrder(ctx context.Context, year int) ([]*record.OrderMonthlyRecord, error)
	GetYearlyOrder(ctx context.Context, year int) ([]*record.OrderYearlyRecord, error)
}

type OrderStatsByMerchantRepository interface {
	GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*record.OrderMonthlyTotalRevenueRecord, error)
	GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*record.OrderYearlyTotalRevenueRecord, error)

	GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*record.OrderMonthlyRecord, error)
	GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*record.OrderYearlyRecord, error)
}

type OrderCommandRepository interface {
	CreateOrder(ctx context.Context, req *requests.CreateOrderRecordRequest) (*record.OrderRecord, error)
	UpdateOrder(ctx context.Context, req *requests.UpdateOrderRecordRequest) (*record.OrderRecord, error)
	TrashedOrder(ctx context.Context, orderID int) (*record.OrderRecord, error)
	RestoreOrder(ctx context.Context, orderID int) (*record.OrderRecord, error)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAllOrder(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}

type OrderQueryRepository interface {
	FindAllOrders(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindById(ctx context.Context, orderID int) (*record.OrderRecord, error)
}
