package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantQueryRepository interface {
	FindById(ctx context.Context, id int) (*db.GetMerchantByIDRow, error)
}

type MerchantBusinessQueryRepository interface {
	FindAllMerchants(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationTrashedRow, error)

	FindById(
		ctx context.Context,
		user_id int,
	) (*db.GetMerchantBusinessInformationRow, error)
}

type MerchantBusinessCommandRepository interface {
	CreateMerchantBusiness(
		ctx context.Context,
		request *requests.CreateMerchantBusinessInformationRequest,
	) (*db.CreateMerchantBusinessInformationRow, error)

	UpdateMerchantBusiness(
		ctx context.Context,
		request *requests.UpdateMerchantBusinessInformationRequest,
	) (*db.UpdateMerchantBusinessInformationRow, error)

	TrashedMerchantBusiness(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantBusinessInformation, error)

	RestoreMerchantBusiness(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantBusinessInformation, error)

	DeleteMerchantBusinessPermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAllMerchantBusiness(ctx context.Context) (bool, error)
	DeleteAllMerchantBusinessPermanent(ctx context.Context) (bool, error)
}
