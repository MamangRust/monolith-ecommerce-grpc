package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
}

type ProductQueryRepository interface {
	FindById(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
}

type ReviewQueryRepository interface {
	FindAllReview(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, error)
	FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, error)
	FindById(ctx context.Context, id int) (*db.GetReviewByIDRow, error)
}

type ReviewCommandRepository interface {
	CreateReview(ctx context.Context, request *requests.CreateReviewRequest) (*db.CreateReviewRow, error)
	UpdateReview(ctx context.Context, request *requests.UpdateReviewRequest) (*db.UpdateReviewRow, error)
	TrashReview(ctx context.Context, shipping_id int) (*db.Review, error)
	RestoreReview(ctx context.Context, category_id int) (*db.Review, error)

	DeleteReviewPermanently(
		ctx context.Context,
		id int,
	) (bool, error)

	RestoreAllReview(ctx context.Context) (bool, error)
	DeleteAllPermanentReview(ctx context.Context) (bool, error)
}
