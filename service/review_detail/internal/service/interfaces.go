package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ReviewDetailQueryService interface {
	FindAllReviews(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsTrashedRow, *int, error)
	FindById(ctx context.Context, user_id int) (*db.GetReviewDetailRow, error)
}

type ReviewDetailCommandService interface {
	CreateReviewDetail(ctx context.Context, request *requests.CreateReviewDetailRequest) (*db.CreateReviewDetailRow, error)
	UpdateReviewDetail(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*db.UpdateReviewDetailRow, error)
	TrashedReviewDetail(ctx context.Context, ReviewDetail_id int) (*db.ReviewDetail, error)
	RestoreReviewDetail(ctx context.Context, ReviewDetail_id int) (*db.ReviewDetail, error)

	DeleteReviewDetailPermanent(
		ctx context.Context,
		review_detail_id int,
	) (bool, error)

	RestoreAllReviewDetail(ctx context.Context) (bool, error)
	DeleteAllReviewDetailPermanent(ctx context.Context) (bool, error)
}
