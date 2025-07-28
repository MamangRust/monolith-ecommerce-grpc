package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type BannerQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, bannerID int) (*response.BannerResponse, *response.ErrorResponse)
}

type BannerCommandService interface {
	CreateBanner(ctx context.Context, req *requests.CreateBannerRequest) (*response.BannerResponse, *response.ErrorResponse)
	UpdateBanner(ctx context.Context, req *requests.UpdateBannerRequest) (*response.BannerResponse, *response.ErrorResponse)
	TrashedBanner(ctx context.Context, bannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse)
	RestoreBanner(ctx context.Context, bannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse)
	DeleteBannerPermanent(ctx context.Context, bannerID int) (bool, *response.ErrorResponse)
	RestoreAllBanner(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllBannerPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
