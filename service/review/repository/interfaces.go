package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
}

type ProductQueryRepository interface {
	FindByID(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
}

type ReviewQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, error)
	FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, error)
	FindActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, error)
	FindByID(ctx context.Context, id int) (*db.GetReviewByIDRow, error)
}

type ReviewCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateReviewRequest) (*db.CreateReviewRow, error)
	Update(ctx context.Context, request *requests.UpdateReviewRequest) (*db.UpdateReviewRow, error)
	Trash(ctx context.Context, review_id int) (*db.Review, error)
	Restore(ctx context.Context, review_id int) (*db.Review, error)

	DeletePermanent(
		ctx context.Context,
		review_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
