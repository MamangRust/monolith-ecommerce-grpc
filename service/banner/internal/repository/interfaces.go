package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type BannerQueryRepository interface {
	FindAllBanners(req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error)
	FindByActive(req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error)
	FindByTrashed(req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error)
	FindById(user_id int) (*record.BannerRecord, error)
}

type BannerCommandRepository interface {
	CreateBanner(request *requests.CreateBannerRequest) (*record.BannerRecord, error)
	UpdateBanner(request *requests.UpdateBannerRequest) (*record.BannerRecord, error)
	TrashedBanner(Banner_id int) (*record.BannerRecord, error)
	RestoreBanner(Banner_id int) (*record.BannerRecord, error)
	DeleteBannerPermanent(Banner_id int) (bool, error)
	RestoreAllBanner() (bool, error)
	DeleteAllBannerPermanent() (bool, error)
}
