package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}

type MerchantDetailQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsRow, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsTrashedRow, error)
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantDetailRow, error)
	FindByIDTrashed(ctx context.Context, user_id int) (*db.GetMerchantDetailTrashedRow, error)
}

type MerchantDetailCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*db.CreateMerchantDetailRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*db.UpdateMerchantDetailRow, error)
	Trash(ctx context.Context, merchant_id int) (*db.MerchantDetail, error)
	Restore(ctx context.Context, merchant_id int) (*db.MerchantDetail, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantSocialLinkCommandRepository interface {
	Create(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*db.CreateMerchantSocialMediaLinkRow, error)
	Update(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*db.UpdateMerchantSocialMediaLinkRow, error)
	Trash(ctx context.Context, socialID int) (bool, error)
	Restore(ctx context.Context, socialID int) (bool, error)
	DeletePermanent(ctx context.Context, socialID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
