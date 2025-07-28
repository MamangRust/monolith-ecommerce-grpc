package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	MerchantQuery         MerchantQueryRepository
	MerchantBusinessQuery MerchantBusinessQueryRepository
	MerchantBusinessCmd   MerchantBusinessCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapper := recordmapper.NewMerchantRecordMapper()
	mapperBusiness := recordmapper.NewMerchantBusinessRecordMapper()

	return &Repositories{
		MerchantQuery:         NewMerchantQueryRepository(DB, mapper),
		MerchantBusinessQuery: NewMerchantBusinessQueryRepository(DB, mapperBusiness),
		MerchantBusinessCmd:   NewMerchantBusinessCommandRepository(DB, mapperBusiness),
	}
}
