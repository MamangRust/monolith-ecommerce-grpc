package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantDetailQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsTrashedRow, *int, error)
	FindById(ctx context.Context, user_id int) (*db.GetMerchantDetailRow, error)
}

type MerchantDetailCommandService interface {
	CreateMerchant(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*db.CreateMerchantDetailRow, error)
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*db.UpdateMerchantDetailRow, error)
	TrashedMerchant(ctx context.Context, merchant_id int) (*db.MerchantDetail, error)
	RestoreMerchant(ctx context.Context, merchant_id int) (*db.MerchantDetail, error)
	DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type MerchantSocialLinkCommandService interface {
	CreateSocialLink(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*db.CreateMerchantSocialMediaLinkRow, error)
	UpdateSocialLink(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*db.UpdateMerchantSocialMediaLinkRow, error)
	TrashSocialLink(ctx context.Context, socialID int) (bool, error)
	RestoreSocialLink(ctx context.Context, socialID int) (bool, error)
	DeletePermanentSocialLink(ctx context.Context, socialID int) (bool, error)
	RestoreAllSocialLink(ctx context.Context) (bool, error)
	DeleteAllPermanentSocialLink(ctx context.Context) (bool, error)
}
