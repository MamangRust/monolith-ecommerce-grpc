package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	MerchantPolicyQuery MerchantPoliciesQueryRepository
	MerchantPolicyCmd   MerchantPoliciesCommandRepository
	MerchantQuery       MerchantQueryRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps Deps) *Repositories {
	mapper := recordmapper.NewMerchantRecordMapper()
	mapperPolicy := recordmapper.NewMerchantPolicyRecordMapper()

	return &Repositories{
		MerchantPolicyQuery: NewMerchantPolicyQueryRepository(deps.DB, deps.Ctx, mapperPolicy),
		MerchantPolicyCmd:   NewMerchantPolicyCommandRepository(deps.DB, deps.Ctx, mapperPolicy),
		MerchantQuery:       NewMerchantQueryRepository(deps.DB, deps.Ctx, mapper),
	}
}
