package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ReviewQueryCache interface {
	GetReviewAllCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, *int, bool)
	SetReviewAllCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewsRow, total *int)

	GetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, *int, bool)
	SetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct, data []*db.GetReviewByProductIdRow, total *int)

	GetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, *int, bool)
	SetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant, data []*db.GetReviewByMerchantIdRow, total *int)

	GetReviewActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, *int, bool)
	SetReviewActiveCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewsActiveRow, total *int)

	GetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, *int, bool)
	SetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewsTrashedRow, total *int)

	GetReviewByIdCache(ctx context.Context, id int) (*db.GetReviewByIDRow, bool)
	SetReviewByIdCache(ctx context.Context, data *db.GetReviewByIDRow)
}

type ReviewCommandCache interface {
	DeleteReviewCache(ctx context.Context, reviewID int)
}
