package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ReviewDetailQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsRow, error)
	FindActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsTrashedRow, error)
	FindByID(ctx context.Context, user_id int) (*db.GetReviewDetailRow, error)
	FindByIDTrashed(ctx context.Context, user_id int) (*db.ReviewDetail, error)
}

type ReviewDetailCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateReviewDetailRequest) (*db.CreateReviewDetailRow, error)
	Update(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*db.UpdateReviewDetailRow, error)
	Trash(ctx context.Context, ReviewDetail_id int) (*db.ReviewDetail, error)
	Restore(ctx context.Context, ReviewDetail_id int) (*db.ReviewDetail, error)

	DeletePermanent(
		ctx context.Context,
		review_detail_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
