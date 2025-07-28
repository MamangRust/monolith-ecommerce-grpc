package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantAwardQueryRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error)
	FindById(ctx context.Context, user_id int) (*record.MerchantAwardRecord, error)
}

type MerchantAwardCommandRepository interface {
	CreateMerchantAward(ctx context.Context, request *requests.CreateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error)
	UpdateMerchantAward(ctx context.Context, request *requests.UpdateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error)
	TrashedMerchantAward(ctx context.Context, merchant_id int) (*record.MerchantAwardRecord, error)
	RestoreMerchantAward(ctx context.Context, merchant_id int) (*record.MerchantAwardRecord, error)
	DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAllMerchantAward(ctx context.Context) (bool, error)
	DeleteAllMerchantAwardPermanent(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*record.MerchantRecord, error)
}
