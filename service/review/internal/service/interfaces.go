package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewQueryService interface {
	FindAllReviews(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse)
	FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse)
	FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse)
}

type ReviewCommandService interface {
	CreateReview(ctx context.Context, req *requests.CreateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse)
	UpdateReview(ctx context.Context, req *requests.UpdateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse)
	TrashedReview(ctx context.Context, reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse)
	RestoreReview(ctx context.Context, reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse)
	DeleteReviewPermanent(ctx context.Context, reviewID int) (bool, *response.ErrorResponse)
	RestoreAllReviews(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllReviewsPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
