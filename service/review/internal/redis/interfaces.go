package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewQueryCache interface {
	GetReviewAllCache(req *requests.FindAllReview) ([]*response.ReviewResponse, *int, bool)
	SetReviewAllCache(req *requests.FindAllReview, data []*response.ReviewResponse, total *int)

	GetReviewByProductCache(req *requests.FindAllReviewByProduct) ([]*response.ReviewsDetailResponse, *int, bool)
	SetReviewByProductCache(req *requests.FindAllReviewByProduct, data []*response.ReviewsDetailResponse, total *int)

	GetReviewByMerchantCache(req *requests.FindAllReviewByMerchant) ([]*response.ReviewsDetailResponse, *int, bool)
	SetReviewByMerchantCache(req *requests.FindAllReviewByMerchant, data []*response.ReviewsDetailResponse, total *int)

	GetReviewActiveCache(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, bool)
	SetReviewActiveCache(req *requests.FindAllReview, data []*response.ReviewResponseDeleteAt, total *int)

	GetReviewTrashedCache(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, bool)
	SetReviewTrashedCache(req *requests.FindAllReview, data []*response.ReviewResponseDeleteAt, total *int)
}
