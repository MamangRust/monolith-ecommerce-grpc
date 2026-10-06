package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
)

type CategoryQueryHandler interface {
	pbcategory.CategoryQueryServiceServer
}

type CategoryCommandHandler interface {
	pbcategory.CategoryCommandServiceServer
}

type CategoryStatsHandler interface {
	pbcategory.CategoryStatsServiceServer
}

type CategoryStatsByIdHandler interface {
	pbcategory.CategoryStatsByIdServiceServer
}

type CategoryStatsByMerchantHandler interface {
	pbcategory.CategoryStatsByMerchantServiceServer
}

type CategoryHandleGrpc interface {
	FindAll(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategory, error)
	FindById(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategory, error)
	FindByActive(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategoryDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategoryDeleteAt, error)
	Create(ctx context.Context, request *pbcategory.CreateCategoryRequest) (*pbcategory.ApiResponseCategory, error)
	Update(ctx context.Context, request *pbcategory.UpdateCategoryRequest) (*pbcategory.ApiResponseCategory, error)
	TrashedCategory(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategoryDeleteAt, error)
	RestoreCategory(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategoryDeleteAt, error)
	DeleteCategoryPermanent(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategoryDelete, error)
	RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pbcategory.ApiResponseCategoryAll, error)
	DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pbcategory.ApiResponseCategoryAll, error)
}
