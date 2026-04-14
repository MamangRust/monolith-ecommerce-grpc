package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ReviewDetailQueryCache interface {
	GetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsRow, *int, bool)
	SetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewDetailsRow, total *int)

	GetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsActiveRow, *int, bool)
	SetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewDetailsActiveRow, total *int)

	GetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsTrashedRow, *int, bool)
	SetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewDetailsTrashedRow, total *int)

	GetCachedReviewDetailCache(ctx context.Context, reviewID int) (*db.GetReviewDetailRow, bool)
	SetCachedReviewDetailCache(ctx context.Context, data *db.GetReviewDetailRow)

	GetCachedReviewDetailTrashedCache(ctx context.Context, reviewID int) (*db.ReviewDetail, bool)
	SetCachedReviewDetailTrashedCache(ctx context.Context, data *db.ReviewDetail)
}

type ReviewDetailCommandCache interface {
	DeleteReviewDetailCache(ctx context.Context, reviewID int)
}
