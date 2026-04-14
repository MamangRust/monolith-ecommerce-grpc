package review_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewQueryCache interface {
	GetReviewAllCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReview, bool)
	SetReviewAllCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReview)

	GetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct) (*response.ApiResponsePaginationReview, bool)
	SetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct, data *response.ApiResponsePaginationReview)

	GetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant) (*response.ApiResponsePaginationReviewsDetail, bool)
	SetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant, data *response.ApiResponsePaginationReviewsDetail)

	GetReviewActiveCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReviewDeleteAt, bool)
	SetReviewActiveCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReviewDeleteAt)

	GetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReviewDeleteAt, bool)
	SetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReviewDeleteAt)

	GetReviewByIdCache(ctx context.Context, id int) (*response.ApiResponseReview, bool)
	SetReviewByIdCache(ctx context.Context, data *response.ApiResponseReview)
}

type ReviewCommandCache interface {
	DeleteReviewCache(ctx context.Context, review_id int)
}
