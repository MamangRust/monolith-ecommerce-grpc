package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantAwardQueryRepository interface {
	FindAllMerchants(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantCertificationsAndAwardsRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantCertificationsAndAwardsActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllMerchant,
	) ([]*db.GetMerchantCertificationsAndAwardsTrashedRow, error)

	FindById(
		ctx context.Context,
		user_id int,
	) (*db.GetMerchantCertificationOrAwardRow, error)
}

type MerchantAwardCommandRepository interface {
	CreateMerchantAward(
		ctx context.Context,
		request *requests.CreateMerchantCertificationOrAwardRequest,
	) (*db.CreateMerchantCertificationOrAwardRow, error)

	UpdateMerchantAward(ctx context.Context, request *requests.UpdateMerchantCertificationOrAwardRequest) (*db.UpdateMerchantCertificationOrAwardRow, error)

	TrashedMerchantAward(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantCertificationsAndAward, error)

	RestoreMerchantAward(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantCertificationsAndAward, error)

	DeleteMerchantPermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAllMerchantAward(ctx context.Context) (bool, error)
	DeleteAllMerchantAwardPermanent(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}
