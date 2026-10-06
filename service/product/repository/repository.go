package repository

import (
	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	categoryadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/category"
	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	Category []adapter.GuardOption
	Merchant []adapter.GuardOption
}

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  CategoryQueryRepository
	MerchantQuery  MerchantQueryRepository
}

func NewRepositories(db *db.Queries,
	categoryQueryClient pbcategory.CategoryQueryServiceClient,
	merchantQueryClient pbmerchant.MerchantQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		ProductQuery:   NewProductQueryRepository(db),
		ProductCommand: NewProductCommandRepository(db),
		CategoryQuery:  categoryadapter.New(categoryQueryClient, g.Category...),
		MerchantQuery:  merchantadapter.New(merchantQueryClient, g.Merchant...),
	}
}
