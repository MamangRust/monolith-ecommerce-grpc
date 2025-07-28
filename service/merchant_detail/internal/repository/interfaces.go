package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*record.MerchantRecord, error)
}

type MerchantDetailQueryRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error)
	FindById(ctx context.Context, user_id int) (*record.MerchantDetailRecord, error)
	FindByIdTrashed(ctx context.Context, user_id int) (*record.MerchantDetailRecord, error)
}

type MerchantDetailCommandRepository interface {
	CreateMerchantDetail(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*record.MerchantDetailRecord, error)
	UpdateMerchantDetail(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*record.MerchantDetailRecord, error)

	TrashedMerchantDetail(ctx context.Context, merchant_detail_id int) (*record.MerchantDetailRecord, error)
	RestoreMerchantDetail(ctx context.Context, merchant_detail_id int) (*record.MerchantDetailRecord, error)
	DeleteMerchantDetailPermanent(ctx context.Context, merchant_detail_id int) (bool, error)
	RestoreAllMerchantDetail(ctx context.Context) (bool, error)
	DeleteAllMerchantDetailPermanent(ctx context.Context) (bool, error)
}

type MerchantSocialLinkCommandRepository interface {
	CreateSocialLink(ctx context.Context, req *requests.CreateMerchantSocialRequest) (bool, error)
	UpdateSocialLink(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (bool, error)
	TrashSocialLink(ctx context.Context, socialID int) (bool, error)
	RestoreSocialLink(ctx context.Context, socialID int) (bool, error)
	DeletePermanentSocialLink(ctx context.Context, socialID int) (bool, error)
	RestoreAllSocialLink(ctx context.Context) (bool, error)
	DeleteAllPermanentSocialLink(ctx context.Context) (bool, error)
}
