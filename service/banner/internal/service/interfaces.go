package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type BannerQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, *int, error)
	FindById(ctx context.Context, bannerID int) (*db.GetBannerRow, error)
}

type BannerCommandService interface {
	CreateBanner(ctx context.Context, req *requests.CreateBannerRequest) (*db.CreateBannerRow, error)
	UpdateBanner(ctx context.Context, req *requests.UpdateBannerRequest) (*db.UpdateBannerRow, error)

	TrashedBanner(ctx context.Context, bannerID int) (*db.Banner, error)
	RestoreBanner(ctx context.Context, bannerID int) (*db.Banner, error)
	DeleteBannerPermanent(ctx context.Context, bannerID int) (bool, error)

	RestoreAllBanner(ctx context.Context) (bool, error)
	DeleteAllBannerPermanent(ctx context.Context) (bool, error)
}
