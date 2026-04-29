package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantAwardQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsActiveRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsTrashedRow, *int, error)
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantCertificationOrAwardRow, error)
}

type MerchantAwardCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantCertificationOrAwardRequest) (*db.CreateMerchantCertificationOrAwardRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantCertificationOrAwardRequest) (*db.UpdateMerchantCertificationOrAwardRow, error)
	Trash(ctx context.Context, merchant_id int) (*db.MerchantCertificationsAndAward, error)
	Restore(ctx context.Context, merchant_id int) (*db.MerchantCertificationsAndAward, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
