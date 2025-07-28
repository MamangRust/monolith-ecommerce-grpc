package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*record.MerchantRecord, error)
}

type MerchantBusinessQueryRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error)
	FindById(ctx context.Context, user_id int) (*record.MerchantBusinessRecord, error)
}

type MerchantBusinessCommandRepository interface {
	CreateMerchantBusiness(ctx context.Context, request *requests.CreateMerchantBusinessInformationRequest) (*record.MerchantBusinessRecord, error)
	UpdateMerchantBusiness(ctx context.Context, request *requests.UpdateMerchantBusinessInformationRequest) (*record.MerchantBusinessRecord, error)
	TrashedMerchantBusiness(ctx context.Context, merchant_id int) (*record.MerchantBusinessRecord, error)
	RestoreMerchantBusiness(ctx context.Context, merchant_id int) (*record.MerchantBusinessRecord, error)
	DeleteMerchantBusinessPermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAllMerchantBusiness(ctx context.Context) (bool, error)
	DeleteAllMerchantBusinessPermanent(ctx context.Context) (bool, error)
}
