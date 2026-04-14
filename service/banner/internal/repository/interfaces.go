package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type BannerQueryRepository interface {
	FindAllBanners(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, error)
	FindById(ctx context.Context, user_id int) (*db.GetBannerRow, error)
}

type BannerCommandRepository interface {
	CreateBanner(ctx context.Context, request *requests.CreateBannerRequest) (*db.CreateBannerRow, error)
	UpdateBanner(ctx context.Context, request *requests.UpdateBannerRequest) (*db.UpdateBannerRow, error)

	TrashedBanner(ctx context.Context, Banner_id int) (*db.Banner, error)
	RestoreBanner(ctx context.Context, Banner_id int) (*db.Banner, error)
	DeleteBannerPermanent(ctx context.Context, banner_id int) (bool, error)

	RestoreAllBanner(ctx context.Context) (bool, error)
	DeleteAllBannerPermanent(ctx context.Context) (bool, error)
}
