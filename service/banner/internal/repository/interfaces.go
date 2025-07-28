package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type BannerQueryRepository interface {
	FindAllBanners(ctx context.Context, req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error)
	FindById(ctx context.Context, user_id int) (*record.BannerRecord, error)
}

type BannerCommandRepository interface {
	CreateBanner(ctx context.Context, request *requests.CreateBannerRequest) (*record.BannerRecord, error)
	UpdateBanner(ctx context.Context, request *requests.UpdateBannerRequest) (*record.BannerRecord, error)
	TrashedBanner(ctx context.Context, banner_id int) (*record.BannerRecord, error)
	RestoreBanner(ctx context.Context, banner_id int) (*record.BannerRecord, error)
	DeleteBannerPermanent(ctx context.Context, banner_id int) (bool, error)
	RestoreAllBanner(ctx context.Context) (bool, error)
	DeleteAllBannerPermanent(ctx context.Context) (bool, error)
}
