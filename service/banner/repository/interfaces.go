package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type BannerQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, error)
	FindActive(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, error)
	FindByID(ctx context.Context, banner_id int) (*db.GetBannerRow, error)
}

type BannerCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateBannerRequest) (*db.CreateBannerRow, error)
	Update(ctx context.Context, request *requests.UpdateBannerRequest) (*db.UpdateBannerRow, error)

	Trash(ctx context.Context, banner_id int) (*db.Banner, error)
	Restore(ctx context.Context, banner_id int) (*db.Banner, error)
	DeletePermanent(ctx context.Context, banner_id int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
