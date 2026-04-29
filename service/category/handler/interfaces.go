package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryQueryHandler interface {
	pb.CategoryQueryServiceServer
}

type CategoryCommandHandler interface {
	pb.CategoryCommandServiceServer
}

type CategoryStatsHandler interface {
	pb.CategoryStatsServiceServer
}

type CategoryStatsByIdHandler interface {
	pb.CategoryStatsByIdServiceServer
}

type CategoryStatsByMerchantHandler interface {
	pb.CategoryStatsByMerchantServiceServer
}

type CategoryHandleGrpc interface {
	FindAll(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategory, error)
	FindById(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategory, error)
	FindByActive(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error)
	Create(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.ApiResponseCategory, error)
	Update(ctx context.Context, request *pb.UpdateCategoryRequest) (*pb.ApiResponseCategory, error)
	TrashedCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error)
	RestoreCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error)
	DeleteCategoryPermanent(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDelete, error)
	RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error)
	DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error)
}
