package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  CategoryQueryRepository
	MerchantQuery  MerchantQueryRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps Deps) *Repositories {
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperCategory := recordmapper.NewCategoryRecordMapper()
	mapperProduct := recordmapper.NewProductRecordMapper()

	return &Repositories{
		ProductQuery:   NewProductQueryRepository(deps.DB, deps.Ctx, mapperProduct),
		ProductCommand: NewProductCommandRepository(deps.DB, deps.Ctx, mapperProduct),
		CategoryQuery:  NewCategoryQueryRepository(deps.DB, deps.Ctx, mapperCategory),
		MerchantQuery:  NewMerchantQueryRepository(deps.DB, deps.Ctx, mapperMerchant),
	}
}
