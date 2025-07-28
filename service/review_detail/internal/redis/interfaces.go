package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewDetailQueryCache interface {
	GetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponse, *int, bool)
	SetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview, data []*response.ReviewDetailsResponse, total *int)

	GetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, bool)
	SetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview, data []*response.ReviewDetailsResponseDeleteAt, total *int)

	GetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, bool)
	SetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*response.ReviewDetailsResponseDeleteAt, total *int)

	GetCachedReviewDetailCache(ctx context.Context, reviewID int) (*response.ReviewDetailsResponse, bool)
	SetCachedReviewDetailCache(ctx context.Context, data *response.ReviewDetailsResponse)
}

type ReviewDetailCommandCache interface {
	DeleteReviewDetailCache(ctx context.Context, reviewID int)
}
