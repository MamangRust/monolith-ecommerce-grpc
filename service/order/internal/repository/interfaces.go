package repository

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type UserQueryRepository interface {
	FindById(user_id int) (*record.UserRecord, error)
}

type ProductQueryRepository interface {
	FindById(user_id int) (*record.ProductRecord, error)
}

type MerchantQueryRepository interface {
	FindById(user_id int) (*record.MerchantRecord, error)
}

type ProductCommandRepository interface {
	UpdateProductCountStock(product_id int, stock int) (*record.ProductRecord, error)
}

type ShippingAddressCommandRepository interface {
	CreateShippingAddress(request *requests.CreateShippingAddressRequest) (*record.ShippingAddressRecord, error)
	UpdateShippingAddress(request *requests.UpdateShippingAddressRequest) (*record.ShippingAddressRecord, error)
}

type OrderItemQueryRepository interface {
	FindOrderItemByOrder(order_id int) ([]*record.OrderItemRecord, error)
	CalculateTotalPrice(order_id int) (*int32, error)
}

type OrderItemCommandRepository interface {
	CreateOrderItem(req *requests.CreateOrderItemRecordRequest) (*record.OrderItemRecord, error)
	UpdateOrderItem(req *requests.UpdateOrderItemRecordRequest) (*record.OrderItemRecord, error)
	TrashedOrderItem(order_id int) (*record.OrderItemRecord, error)
	RestoreOrderItem(order_id int) (*record.OrderItemRecord, error)
	DeleteOrderItemPermanent(order_id int) (bool, error)
	RestoreAllOrderItem() (bool, error)
	DeleteAllOrderPermanent() (bool, error)
}

type OrderStatsRepository interface {
	GetMonthlyTotalRevenue(req *requests.MonthTotalRevenue) ([]*record.OrderMonthlyTotalRevenueRecord, error)
	GetYearlyTotalRevenue(year int) ([]*record.OrderYearlyTotalRevenueRecord, error)

	GetMonthlyOrder(year int) ([]*record.OrderMonthlyRecord, error)
	GetYearlyOrder(year int) ([]*record.OrderYearlyRecord, error)
}

type OrderStatsByMerchantRepository interface {
	GetMonthlyTotalRevenueByMerchant(req *requests.MonthTotalRevenueMerchant) ([]*record.OrderMonthlyTotalRevenueRecord, error)
	GetYearlyTotalRevenueByMerchant(req *requests.YearTotalRevenueMerchant) ([]*record.OrderYearlyTotalRevenueRecord, error)

	GetMonthlyOrderByMerchant(req *requests.MonthOrderMerchant) ([]*record.OrderMonthlyRecord, error)
	GetYearlyOrderByMerchant(req *requests.YearOrderMerchant) ([]*record.OrderYearlyRecord, error)
}

type OrderCommandRepository interface {
	CreateOrder(req *requests.CreateOrderRecordRequest) (*record.OrderRecord, error)
	UpdateOrder(req *requests.UpdateOrderRecordRequest) (*record.OrderRecord, error)
	TrashedOrder(order_id int) (*record.OrderRecord, error)
	RestoreOrder(order_id int) (*record.OrderRecord, error)
	DeleteOrderPermanent(order_id int) (bool, error)
	RestoreAllOrder() (bool, error)
	DeleteAllOrderPermanent() (bool, error)
}

type OrderQueryRepository interface {
	FindAllOrders(req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindByActive(req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindByTrashed(req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error)
	FindById(order_id int) (*record.OrderRecord, error)
}
