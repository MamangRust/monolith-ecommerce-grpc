package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderQueryHandler interface {
	pb.OrderQueryServiceServer
}

type OrderCommandHandler interface {
	pb.OrderCommandServiceServer
}

type OrderStatsHandler interface {
	pb.OrderQueryServiceServer
}

type OrderStatsByMerchantHandler interface {
	pb.OrderQueryServiceServer
}

type OrderHandleGrpc interface {
	FindMonthlyTotalRevenue(ctx context.Context, request *pb.FindYearMonthTotalRevenue) (*pb.ApiResponseOrderMonthlyTotalRevenue, error)
	FindYearlyTotalRevenue(ctx context.Context, request *pb.FindYearTotalRevenue) (*pb.ApiResponseOrderYearlyTotalRevenue, error)

	FindMonthlyTotalRevenueByMerchant(ctx context.Context, request *pb.FindYearMonthTotalRevenueByMerchant) (*pb.ApiResponseOrderMonthlyTotalRevenue, error)
	FindYearlyTotalRevenueByMerchant(ctx context.Context, request *pb.FindYearTotalRevenueByMerchant) (*pb.ApiResponseOrderYearlyTotalRevenue, error)

	FindMonthlyRevenue(ctx context.Context, request *pb.FindYearOrder) (*pb.ApiResponseOrderMonthly, error)
	FindYearlyRevenue(ctx context.Context, request *pb.FindYearOrder) (*pb.ApiResponseOrderYearly, error)

	FindMonthlyRevenueByMerchant(ctx context.Context, request *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderMonthly, error)
	FindYearlyRevenueByMerchant(ctx context.Context, request *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderYearly, error)

	FindAll(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrder, error)
	FindById(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrder, error)

	FindByActive(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error)

	Create(ctx context.Context, request *pb.CreateOrderRequest) (*pb.ApiResponseOrder, error)
	Update(ctx context.Context, request *pb.UpdateOrderRequest) (*pb.ApiResponseOrder, error)
	TrashedOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error)
	RestoreOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error)
	DeleteOrderPermanent(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDelete, error)
	RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error)
	DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error)
}
