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

func NewRepositories(DB *db.Queries) *Repositories {
	merchantMapper := recordmapper.NewMerchantRecordMapper()
	merchantAwardMapper := recordmapper.NewMerchantAwardRecordMapper()

	return &Repositories{
		MerchantAwardQuery:   NewMerchantAwardQueryRepository(DB, merchantAwardMapper),
		MerchantAwardCommand: NewMerchantAwardCommandRepository(DB, merchantAwardMapper),
		MerchantQuery:        NewMerchantQueryRepository(DB, merchantMapper),
	}
}
