package reviewdetail_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewDetailQueryCache interface {
	GetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReviewDetails, bool)
	SetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReviewDetails)

	GetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReviewDetailsDeleteAt, bool)
	SetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReviewDetailsDeleteAt)

	GetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReviewDetailsDeleteAt, bool)
	SetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReviewDetailsDeleteAt)

	GetCachedReviewDetailCache(ctx context.Context, reviewID int) (*response.ApiResponseReviewDetail, bool)
	SetCachedReviewDetailCache(ctx context.Context, data *response.ApiResponseReviewDetail)

	GetCachedReviewDetailTrashedCache(ctx context.Context, reviewID int) (*response.ApiResponseReviewDetailDeleteAt, bool)
	SetCachedReviewDetailTrashedCache(ctx context.Context, data *response.ApiResponseReviewDetailDeleteAt)
}

type ReviewDetailCommandCache interface {
	DeleteReviewDetailCache(ctx context.Context, reviewID int)
}
