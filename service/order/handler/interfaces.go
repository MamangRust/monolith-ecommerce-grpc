package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pborder "github.com/MamangRust/monolith-ecommerce-pb/order"
)

type OrderQueryHandler interface {
	pborder.OrderQueryServiceServer
}

type OrderCommandHandler interface {
	pborder.OrderCommandServiceServer
}

type OrderStatsHandler interface {
	pborder.OrderStatsServiceServer
}

type OrderStatsByMerchantHandler interface {
	pborder.OrderStatsByMerchantServiceServer
}

type OrderHandleGrpc interface {
	FindMonthlyTotalRevenue(ctx context.Context, request *pborder.FindYearMonthTotalRevenue) (*pborder.ApiResponseOrderMonthlyTotalRevenue, error)
	FindYearlyTotalRevenue(ctx context.Context, request *pborder.FindYearTotalRevenue) (*pborder.ApiResponseOrderYearlyTotalRevenue, error)

	FindMonthlyTotalRevenueByMerchant(ctx context.Context, request *pborder.FindYearMonthTotalRevenueByMerchant) (*pborder.ApiResponseOrderMonthlyTotalRevenue, error)
	FindYearlyTotalRevenueByMerchant(ctx context.Context, request *pborder.FindYearTotalRevenueByMerchant) (*pborder.ApiResponseOrderYearlyTotalRevenue, error)

	FindMonthlyRevenue(ctx context.Context, request *pborder.FindYearOrder) (*pborder.ApiResponseOrderMonthly, error)
	FindYearlyRevenue(ctx context.Context, request *pborder.FindYearOrder) (*pborder.ApiResponseOrderYearly, error)

	FindMonthlyRevenueByMerchant(ctx context.Context, request *pborder.FindYearOrderByMerchant) (*pborder.ApiResponseOrderMonthly, error)
	FindYearlyRevenueByMerchant(ctx context.Context, request *pborder.FindYearOrderByMerchant) (*pborder.ApiResponseOrderYearly, error)

	FindAll(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrder, error)
	FindById(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrder, error)

	FindByActive(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrderDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrderDeleteAt, error)

	Create(ctx context.Context, request *pborder.CreateOrderRequest) (*pborder.ApiResponseOrder, error)
	Update(ctx context.Context, request *pborder.UpdateOrderRequest) (*pborder.ApiResponseOrder, error)
	TrashedOrder(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDeleteAt, error)
	RestoreOrder(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDeleteAt, error)
	DeleteOrderPermanent(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDelete, error)
	RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pborder.ApiResponseOrderAll, error)
	DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pborder.ApiResponseOrderAll, error)
}
