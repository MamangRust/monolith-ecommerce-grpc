package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewDetailQueryService interface {
	FindAll(req *requests.FindAllReview) ([]*response.ReviewDetailsResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(review_id int) (*response.ReviewDetailsResponse, *response.ErrorResponse)
}

type ReviewDetailCommandService interface {
	CreateReviewDetail(req *requests.CreateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse)
	UpdateReviewDetail(req *requests.UpdateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse)
	TrashedReviewDetail(review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse)
	RestoreReviewDetail(review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse)
	DeleteReviewDetailPermanent(review_id int) (bool, *response.ErrorResponse)
	RestoreAllReviewDetail() (bool, *response.ErrorResponse)
	DeleteAllReviewDetailPermanent() (bool, *response.ErrorResponse)
}
