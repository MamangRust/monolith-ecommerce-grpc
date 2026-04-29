package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type Repositories struct {
	MerchantPoliciesQuery   MerchantPoliciesQueryRepository
	MerchantPoliciesCommand MerchantPoliciesCommandRepository
	MerchantQuery           MerchantQueryRepository
}

func NewRepositories(DB *db.Queries, merchantQuery pb.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantPoliciesQuery:   NewMerchantPolicyQueryRepository(DB),
		MerchantPoliciesCommand: NewMerchantPolicyCommandRepository(DB),
		MerchantQuery:           NewMerchantQueryRepository(merchantQuery),
	}
}
