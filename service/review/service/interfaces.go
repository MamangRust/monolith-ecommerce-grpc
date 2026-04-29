package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ReviewQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, *int, error)
	FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, *int, error)
	FindByID(ctx context.Context, id int) (*db.GetReviewByIDRow, error)
}

type ReviewCommandService interface {
	Create(ctx context.Context, request *requests.CreateReviewRequest) (*db.CreateReviewRow, error)
	Update(ctx context.Context, request *requests.UpdateReviewRequest) (*db.UpdateReviewRow, error)
	Trash(ctx context.Context, shipping_id int) (*db.Review, error)
	Restore(ctx context.Context, category_id int) (*db.Review, error)

	DeletePermanent(
		ctx context.Context,
		id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
