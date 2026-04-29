package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type BannerQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, *int, error)
	FindByID(ctx context.Context, bannerID int) (*db.GetBannerRow, error)
}

type BannerCommandService interface {
	Create(ctx context.Context, req *requests.CreateBannerRequest) (*db.CreateBannerRow, error)
	Update(ctx context.Context, req *requests.UpdateBannerRequest) (*db.UpdateBannerRow, error)

	Trash(ctx context.Context, bannerID int) (*db.Banner, error)
	Restore(ctx context.Context, bannerID int) (*db.Banner, error)
	DeletePermanent(ctx context.Context, bannerID int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
