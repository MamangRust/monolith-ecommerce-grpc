package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	CartQuery    CartQueryRepository
	CartCommand  CartCommandRepository
	UserQuery    UserQueryRepository
	ProductQuery ProductQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapperCart := recordmapper.NewCartRecordMapper()
	mapperUser := recordmapper.NewUserRecordMapper()
	mapperProduct := recordmapper.NewProductRecordMapper()

	return &Repositories{
		CartQuery:    NewCartQueryRepository(DB, mapperCart),
		CartCommand:  NewCartCommandRepository(DB, mapperCart),
		UserQuery:    NewUserQueryRepository(DB, mapperUser),
		ProductQuery: NewProductQueryRepository(DB, mapperProduct),
	}
}
