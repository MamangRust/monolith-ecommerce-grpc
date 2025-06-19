package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewDetailQueryCache interface {
	GetReviewDetailAllCache(req *requests.FindAllReview) ([]*response.ReviewDetailsResponse, *int, bool)
	SetReviewDetailAllCache(req *requests.FindAllReview, data []*response.ReviewDetailsResponse, total *int)

	GetRevieDetailActiveCache(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, bool)
	SetReviewDetailActiveCache(req *requests.FindAllReview, data []*response.ReviewDetailsResponseDeleteAt, total *int)

	GetReviewDetailTrashedCache(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, bool)
	SetReviewDetailTrashedCache(req *requests.FindAllReview, data []*response.ReviewDetailsResponseDeleteAt, total *int)

	GetCachedReviewDetailCache(review_id int) (*response.ReviewDetailsResponse, bool)
	SetCachedReviewDetailCache(data *response.ReviewDetailsResponse)
}

type ReviewDetailCommandCache interface {
	DeleteReviewDetailCache(review_id int)
}
