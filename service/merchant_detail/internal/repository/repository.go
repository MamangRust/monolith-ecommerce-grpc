package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	MerchantQuery             MerchantQueryRepository
	MerchantDetailQuery       MerchantDetailQueryRepository
	MerchantDetailCommand     MerchantDetailCommandRepository
	MerchantSocialLinkCommand MerchantSocialLinkCommandRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps *Deps) *Repositories {
	mapper := recordmapper.NewMerchantRecordMapper()
	mapperDetail := recordmapper.NewMerchantDetailRecordMapper()

	return &Repositories{
		MerchantQuery:             NewMerchantQueryRepository(deps.DB, deps.Ctx, mapper),
		MerchantDetailQuery:       NewMerchantDetailQueryRepository(deps.DB, deps.Ctx, mapperDetail),
		MerchantDetailCommand:     NewMerchantDetailCommandRepository(deps.DB, deps.Ctx, mapperDetail),
		MerchantSocialLinkCommand: NewMerchantSocialLinkCommandRepository(deps.DB, deps.Ctx),
	}
}
