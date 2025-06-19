package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderStatsCache interface {
	GetMonthlyTotalRevenueCache(req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, bool)
	SetMonthlyTotalRevenueCache(req *requests.MonthTotalRevenue, data []*response.OrderMonthlyTotalRevenueResponse)

	GetYearlyTotalRevenueCache(year int) ([]*response.OrderYearlyTotalRevenueResponse, bool)
	SetYearlyTotalRevenueCache(year int, data []*response.OrderYearlyTotalRevenueResponse)

	GetMonthlyOrderCache(year int) ([]*response.OrderMonthlyResponse, bool)
	SetMonthlyOrderCache(year int, data []*response.OrderMonthlyResponse)

	GetYearlyOrderCache(year int) ([]*response.OrderYearlyResponse, bool)
	SetYearlyOrderCache(year int, data []*response.OrderYearlyResponse)
}

type OrderStatsByMerchantCache interface {
	GetMonthlyTotalRevenueByMerchantCache(req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, bool)
	SetMonthlyTotalRevenueByMerchantCache(req *requests.MonthTotalRevenueMerchant, data []*response.OrderMonthlyTotalRevenueResponse)

	GetYearlyTotalRevenueByMerchantCache(req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, bool)
	SetYearlyTotalRevenueByMerchantCache(req *requests.YearTotalRevenueMerchant, data []*response.OrderYearlyTotalRevenueResponse)

	GetMonthlyOrderByMerchantCache(req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, bool)
	SetMonthlyOrderByMerchantCache(req *requests.MonthOrderMerchant, data []*response.OrderMonthlyResponse)

	GetYearlyOrderByMerchantCache(req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, bool)
	SetYearlyOrderByMerchantCache(req *requests.YearOrderMerchant, data []*response.OrderYearlyResponse)
}

type OrderQueryCache interface {
	GetOrderAllCache(req *requests.FindAllOrder) ([]*response.OrderResponse, *int, bool)
	SetOrderAllCache(req *requests.FindAllOrder, data []*response.OrderResponse, total *int)

	GetOrderActiveCache(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, bool)
	SetOrderActiveCache(req *requests.FindAllOrder, data []*response.OrderResponseDeleteAt, total *int)

	GetOrderTrashedCache(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, bool)
	SetOrderTrashedCache(req *requests.FindAllOrder, data []*response.OrderResponseDeleteAt, total *int)

	GetCachedOrderCache(order_id int) (*response.OrderResponse, bool)
	SetCachedOrderCache(data *response.OrderResponse)
}

type OrderCommandCache interface {
	DeleteOrderCache(order_id int)
}
