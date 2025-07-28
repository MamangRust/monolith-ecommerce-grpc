package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	BannerQuery   BannerQueryRepository
	BannerCommand BannerCommandRepository
}

func NewRepositories(db *db.Queries) *Repositories {
	bannerMapper := recordmapper.NewBannerRecordMapper()

	return &Repositories{
		BannerQuery:   NewBannerQueryRepository(db, bannerMapper),
		BannerCommand: NewBannerCommandRepository(db, bannerMapper),
	}
}
