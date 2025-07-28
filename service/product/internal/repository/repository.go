package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  CategoryQueryRepository
	MerchantQuery  MerchantQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperCategory := recordmapper.NewCategoryRecordMapper()
	mapperProduct := recordmapper.NewProductRecordMapper()

	return &Repositories{
		ProductQuery:   NewProductQueryRepository(DB, mapperProduct),
		ProductCommand: NewProductCommandRepository(DB, mapperProduct),
		CategoryQuery:  NewCategoryQueryRepository(DB, mapperCategory),
		MerchantQuery:  NewMerchantQueryRepository(DB, mapperMerchant),
	}
}
