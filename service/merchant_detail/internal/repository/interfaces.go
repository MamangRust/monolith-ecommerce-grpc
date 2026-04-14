package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}

type MerchantDetailQueryRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsTrashedRow, error)
	FindById(ctx context.Context, user_id int) (*db.GetMerchantDetailRow, error)
	FindByIdTrashed(ctx context.Context, user_id int) (*db.GetMerchantDetailTrashedRow, error)
}

type MerchantDetailCommandRepository interface {
	CreateMerchantDetail(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*db.CreateMerchantDetailRow, error)
	UpdateMerchantDetail(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*db.UpdateMerchantDetailRow, error)
	TrashedMerchantDetail(ctx context.Context, merchant_id int) (*db.MerchantDetail, error)
	RestoreMerchantDetail(ctx context.Context, merchant_id int) (*db.MerchantDetail, error)
	DeleteMerchantDetailPermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAllMerchantDetail(ctx context.Context) (bool, error)
	DeleteAllMerchantDetailPermanent(ctx context.Context) (bool, error)
}

type MerchantSocialLinkCommandRepository interface {
	CreateSocialLink(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*db.CreateMerchantSocialMediaLinkRow, error)
	UpdateSocialLink(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*db.UpdateMerchantSocialMediaLinkRow, error)
	TrashSocialLink(ctx context.Context, socialID int) (bool, error)
	RestoreSocialLink(ctx context.Context, socialID int) (bool, error)
	DeletePermanentSocialLink(ctx context.Context, socialID int) (bool, error)
	RestoreAllSocialLink(ctx context.Context) (bool, error)
	DeleteAllPermanentSocialLink(ctx context.Context) (bool, error)
}
