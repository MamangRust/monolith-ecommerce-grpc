package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewQueryCache interface {
	GetReviewAllCache(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewResponse, *int, bool)
	SetReviewAllCache(ctx context.Context, req *requests.FindAllReview, data []*response.ReviewResponse, total *int)

	GetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*response.ReviewsDetailResponse, *int, bool)
	SetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct, data []*response.ReviewsDetailResponse, total *int)

	GetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*response.ReviewsDetailResponse, *int, bool)
	SetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant, data []*response.ReviewsDetailResponse, total *int)

	GetReviewActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, bool)
	SetReviewActiveCache(ctx context.Context, req *requests.FindAllReview, data []*response.ReviewResponseDeleteAt, total *int)

	GetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, bool)
	SetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*response.ReviewResponseDeleteAt, total *int)
}
