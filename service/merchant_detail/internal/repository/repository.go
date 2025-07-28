package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	MerchantQuery             MerchantQueryRepository
	MerchantDetailQuery       MerchantDetailQueryRepository
	MerchantDetailCommand     MerchantDetailCommandRepository
	MerchantSocialLinkCommand MerchantSocialLinkCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapper := recordmapper.NewMerchantRecordMapper()
	mapperDetail := recordmapper.NewMerchantDetailRecordMapper()

	return &Repositories{
		MerchantQuery:             NewMerchantQueryRepository(DB, mapper),
		MerchantDetailQuery:       NewMerchantDetailQueryRepository(DB, mapperDetail),
		MerchantDetailCommand:     NewMerchantDetailCommandRepository(DB, mapperDetail),
		MerchantSocialLinkCommand: NewMerchantSocialLinkCommandRepository(DB),
	}
}
