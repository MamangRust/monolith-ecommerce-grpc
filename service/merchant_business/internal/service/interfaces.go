package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantBusinessQueryService interface {
	FindAllMerchants(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantsBusinessInformationTrashedRow, *int, error)

	FindById(
		ctx context.Context,
		user_id int,
	) (*db.GetMerchantBusinessInformationRow, error)
}

type MerchantBusinessCommandService interface {
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
