package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

type Repositories struct {
	BannerQuery   BannerQueryRepository
	BannerCommand BannerCommandRepository
}

func NewRepositories(db *db.Queries) *Repositories {

	return &Repositories{
		BannerQuery:   NewBannerQueryRepository(db),
		BannerCommand: NewBannerCommandRepository(db),
	}
}
