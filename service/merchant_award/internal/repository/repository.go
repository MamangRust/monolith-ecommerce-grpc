package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
	"golang.org/x/net/context"
)

type Repositories struct {
	MerchantAwardQuery   MerchantAwardQueryRepository
	MerchantAwardCommand MerchantAwardCommandRepository
	MerchantQuery        MerchantQueryRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps *Deps) *Repositories {
	merchantMapper := recordmapper.NewMerchantRecordMapper()
	merchantAwardMapper := recordmapper.NewMerchantAwardRecordMapper()

	return &Repositories{
		MerchantAwardQuery:   NewMerchantAwardQueryRepository(deps.DB, deps.Ctx, merchantAwardMapper),
		MerchantAwardCommand: NewMerchantAwardCommandRepository(deps.DB, deps.Ctx, merchantAwardMapper),
		MerchantQuery:        NewMerchantQueryRepository(deps.DB, deps.Ctx, merchantMapper),
	}
}
