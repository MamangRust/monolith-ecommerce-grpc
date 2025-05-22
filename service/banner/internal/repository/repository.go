package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	BannerQuery   BannerQueryRepository
	BannerCommand BannerCommandRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps Deps) *Repositories {
	bannerMapper := recordmapper.NewBannerRecordMapper()

	return &Repositories{
		BannerQuery:   NewBannerQueryRepository(deps.DB, deps.Ctx, bannerMapper),
		BannerCommand: NewBannerCommandRepository(deps.DB, deps.Ctx, bannerMapper),
	}
}
