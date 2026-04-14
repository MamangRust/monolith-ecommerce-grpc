package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

type Repositories struct {
	MerchantQuery         MerchantQueryRepository
	MerchantBusinessQuery MerchantBusinessQueryRepository
	MerchantBusinessCommand MerchantBusinessCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		MerchantQuery:         NewMerchantQueryRepository(DB),
		MerchantBusinessQuery: NewMerchantBusinessQueryRepository(DB),
		MerchantBusinessCommand: NewMerchantBusinessCommandRepository(DB),
	}
}
