package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantCommand         MerchantCommandRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	UserQuery               UserQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapper := recordmapper.NewMerchantRecordMapper()
	mapperDocument := recordmapper.NewMerchantDocumentRecordMapper()
	mapperUser := recordmapper.NewUserRecordMapper()

	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(DB, mapper),
		MerchantCommand:         NewMerchantCommandRepository(DB, mapper),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(DB, mapperDocument),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(DB, mapperDocument),
		UserQuery:               NewUserQueryRepository(DB, mapperUser),
	}
}
