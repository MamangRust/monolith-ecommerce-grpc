package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

type Repositories struct {
	MerchantPoliciesQuery   MerchantPoliciesQueryRepository
	MerchantPoliciesCommand MerchantPoliciesCommandRepository
	MerchantQuery           MerchantQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		MerchantPoliciesQuery:   NewMerchantPolicyQueryRepository(DB),
		MerchantPoliciesCommand: NewMerchantPolicyCommandRepository(DB),
		MerchantQuery:           NewMerchantQueryRepository(DB),
	}
}
