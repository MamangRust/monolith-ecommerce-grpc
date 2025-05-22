package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	CartQuery    CartQueryRepository
	CartCommand  CartCommandRepository
	UserQuery    UserQueryRepository
	ProductQuery ProductQueryRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps Deps) *Repositories {
	mapperCart := recordmapper.NewCartRecordMapper()
	mapperUser := recordmapper.NewUserRecordMapper()
	mapperProduct := recordmapper.NewProductRecordMapper()

	return &Repositories{
		CartQuery:    NewCartQueryRepository(deps.DB, deps.Ctx, mapperCart),
		CartCommand:  NewCartCommandRepository(deps.DB, deps.Ctx, mapperCart),
		UserQuery:    NewUserQueryRepository(deps.DB, deps.Ctx, mapperUser),
		ProductQuery: NewProductQueryRepository(deps.DB, deps.Ctx, mapperProduct),
	}
}
