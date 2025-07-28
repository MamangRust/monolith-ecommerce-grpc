package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ReviewDetailQueryRepository interface {
	FindAllReviews(ctx context.Context, req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error)
	FindById(ctx context.Context, userID int) (*record.ReviewDetailRecord, error)
	FindByIdTrashed(ctx context.Context, userID int) (*record.ReviewDetailRecord, error)
}

type ReviewDetailCommandRepository interface {
	CreateReviewDetail(ctx context.Context, request *requests.CreateReviewDetailRequest) (*record.ReviewDetailRecord, error)
	UpdateReviewDetail(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*record.ReviewDetailRecord, error)
	TrashedReviewDetail(ctx context.Context, reviewDetailID int) (*record.ReviewDetailRecord, error)
	RestoreReviewDetail(ctx context.Context, reviewDetailID int) (*record.ReviewDetailRecord, error)
	DeleteReviewDetailPermanent(ctx context.Context, reviewDetailID int) (bool, error)
	RestoreAllReviewDetail(ctx context.Context) (bool, error)
	DeleteAllReviewDetailPermanent(ctx context.Context) (bool, error)
}
