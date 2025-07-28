package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	MerchantPolicyQuery MerchantPoliciesQueryRepository
	MerchantPolicyCmd   MerchantPoliciesCommandRepository
	MerchantQuery       MerchantQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapper := recordmapper.NewMerchantRecordMapper()
	mapperPolicy := recordmapper.NewMerchantPolicyRecordMapper()

	return &Repositories{
		MerchantPolicyQuery: NewMerchantPolicyQueryRepository(DB, mapperPolicy),
		MerchantPolicyCmd:   NewMerchantPolicyCommandRepository(DB, mapperPolicy),
		MerchantQuery:       NewMerchantQueryRepository(DB, mapper),
	}
}
