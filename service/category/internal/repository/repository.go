package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	CategoryQuery           CategoryQueryRepository
	CategoryCommand         CategoryCommandRepository
	CategoryStats           CategoryStatsRepository
	CategoryStatsById       CategoryStatsByIdRepository
	CategoryStatsByMerchant CategoryStatsByMerchantRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps *Deps) *Repositories {
	categoryMapper := recordmapper.NewCategoryRecordMapper()

	return &Repositories{
		CategoryQuery:           NewCategoryQueryRepository(deps.DB, deps.Ctx, categoryMapper),
		CategoryCommand:         NewCategoryCommandRepository(deps.DB, deps.Ctx, categoryMapper),
		CategoryStats:           NewCategoryStatsRepository(deps.DB, deps.Ctx, categoryMapper),
		CategoryStatsById:       NewCategoryStatsByIdRepository(deps.DB, deps.Ctx, categoryMapper),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantRepository(deps.DB, deps.Ctx, categoryMapper),
	}
}
