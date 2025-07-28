package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewDetailQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, reviewID int) (*response.ReviewDetailsResponse, *response.ErrorResponse)
}

type ReviewDetailCommandService interface {
	CreateReviewDetail(ctx context.Context, req *requests.CreateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse)
	UpdateReviewDetail(ctx context.Context, req *requests.UpdateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse)
	TrashedReviewDetail(ctx context.Context, reviewID int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse)
	RestoreReviewDetail(ctx context.Context, reviewID int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse)
	DeleteReviewDetailPermanent(ctx context.Context, reviewID int) (bool, *response.ErrorResponse)
	RestoreAllReviewDetail(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllReviewDetailPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
