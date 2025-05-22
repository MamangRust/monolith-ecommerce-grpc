package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewQueryService interface {
	FindAllReviews(req *requests.FindAllReview) ([]*response.ReviewResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse)
	FindByProduct(req *requests.FindAllReviewByProduct) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse)
	FindByMerchant(req *requests.FindAllReviewByMerchant) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse)
}

type ReviewCommandService interface {
	CreateReview(req *requests.CreateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse)
	UpdateReview(req *requests.UpdateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse)
	TrashedReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse)
	RestoreReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse)
	DeleteReviewPermanent(reviewID int) (bool, *response.ErrorResponse)
	RestoreAllReviews() (bool, *response.ErrorResponse)
	DeleteAllReviewsPermanent() (bool, *response.ErrorResponse)
}
