package repository

import (
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	Merchant []adapter.GuardOption
}

type Repositories struct {
	MerchantAwardQuery   MerchantAwardQueryRepository
	MerchantAwardCommand MerchantAwardCommandRepository
	MerchantQuery        MerchantQueryRepository
}

func NewRepositories(db *db.Queries, merchantQueryClient pbmerchant.MerchantQueryServiceClient, guards ...GuardOptions) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		MerchantAwardQuery:   NewMerchantAwardQueryRepository(db),
		MerchantAwardCommand: NewMerchantAwardCommandRepository(db),
		MerchantQuery:        merchantadapter.New(merchantQueryClient, g.Merchant...),
	}
}
