package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type BannerQueryService interface {
	FindAll(req *requests.FindAllBanner) ([]*response.BannerResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(BannerID int) (*response.BannerResponse, *response.ErrorResponse)
}

type BannerCommandService interface {
	CreateBanner(req *requests.CreateBannerRequest) (*response.BannerResponse, *response.ErrorResponse)
	UpdateBanner(req *requests.UpdateBannerRequest) (*response.BannerResponse, *response.ErrorResponse)
	TrashedBanner(BannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse)
	RestoreBanner(BannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse)
	DeleteBannerPermanent(BannerID int) (bool, *response.ErrorResponse)
	RestoreAllBanner() (bool, *response.ErrorResponse)
	DeleteAllBannerPermanent() (bool, *response.ErrorResponse)
}
