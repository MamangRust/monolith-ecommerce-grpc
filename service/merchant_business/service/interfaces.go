package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantBusinessQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationActiveRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationTrashedRow, *int, error)
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantBusinessInformationRow, error)
}

type MerchantBusinessCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantBusinessInformationRequest) (*db.CreateMerchantBusinessInformationRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantBusinessInformationRequest) (*db.UpdateMerchantBusinessInformationRow, error)
	Trash(ctx context.Context, merchant_id int) (*db.MerchantBusinessInformation, error)
	Restore(ctx context.Context, merchant_id int) (*db.MerchantBusinessInformation, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
