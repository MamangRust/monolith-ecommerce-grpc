package handler

import (
	pbbanner "github.com/MamangRust/monolith-ecommerce-pb/banner"
)

type BannerQueryHandler interface {
	pbbanner.BannerQueryServiceServer
}

type BannerCommandHandler interface {
	pbbanner.BannerCommandServiceServer
}
