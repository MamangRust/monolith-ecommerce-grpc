package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	MerchantQuery         MerchantQueryRepository
	MerchantBusinessQuery MerchantBusinessQueryRepository
	MerchantBusinessCmd   MerchantBusinessCommandRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps Deps) *Repositories {
	mapper := recordmapper.NewMerchantRecordMapper()
	mapperBusiness := recordmapper.NewMerchantBusinessRecordMapper()

	return &Repositories{
		MerchantQuery:         NewMerchantQueryRepository(deps.DB, deps.Ctx, mapper),
		MerchantBusinessQuery: NewMerchantBusinessQueryRepository(deps.DB, deps.Ctx, mapperBusiness),
		MerchantBusinessCmd:   NewMerchantBusinessCommandRepository(deps.DB, deps.Ctx, mapperBusiness),
	}
}
