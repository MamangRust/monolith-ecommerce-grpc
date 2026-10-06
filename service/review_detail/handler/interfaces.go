package handler

import (
	"context"

	pbreview "github.com/MamangRust/monolith-ecommerce-pb/review"
	pbreview_detail "github.com/MamangRust/monolith-ecommerce-pb/review_detail"
)

type ReviewDetailQueryHandler interface {
	pbreview_detail.ReviewDetailQueryServiceServer
}

type ReviewDetailCommandHandler interface {
	pbreview_detail.ReviewDetailCommandServiceServer
}

type ReviewDetailHandleGrpc interface {
	FindAll(ctx context.Context, request *pbreview.FindAllReviewRequest) (*pbreview_detail.ApiResponsePaginationReviewDetails, error)
	FindById(ctx context.Context, request *pbreview_detail.FindByIdReviewDetailRequest) (*pbreview_detail.ApiResponseReviewDetail, error)
	FindByActive(ctx context.Context, request *pbreview.FindAllReviewRequest) (*pbreview_detail.ApiResponsePaginationReviewDetailsDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pbreview.FindAllReviewRequest) (*pbreview_detail.ApiResponsePaginationReviewDetailsDeleteAt, error)
	Create(ctx context.Context, request *pbreview_detail.CreateReviewDetailRequest) (*pbreview_detail.ApiResponseReviewDetail, error)
	Update(ctx context.Context, request *pbreview_detail.UpdateReviewDetailRequest) (*pbreview_detail.ApiResponseReviewDetail, error)
	TrashedReviewDetail(ctx context.Context, request *pbreview_detail.FindByIdReviewDetailRequest) (*pbreview_detail.ApiResponseReviewDetailDeleteAt, error)
	RestoreReviewDetail(ctx context.Context, request *pbreview_detail.FindByIdReviewDetailRequest) (*pbreview_detail.ApiResponseReviewDetailDeleteAt, error)
	DeleteReviewDetailPermanent(ctx context.Context, request *pbreview_detail.FindByIdReviewDetailRequest) (*pbreview.ApiResponseReviewDelete, error)
}
