package repository

import (
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	productadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/product"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	User    []adapter.GuardOption
	Product []adapter.GuardOption
}

type Repositories struct {
	ProductQuery  ProductQueryRepository
	ReviewQuery   ReviewQueryRepository
	UserQuery     UserQueryRepository
	ReviewCommand ReviewCommandRepository
}

func NewRepositories(queries *db.Queries,
	userQueryClient pbuser.UserQueryServiceClient,
	productQueryClient pbproduct.ProductQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		ProductQuery:  productadapter.NewQueryAdapter(productQueryClient, g.Product...),
		ReviewQuery:   NewReviewQueryRepository(queries),
		UserQuery:     useradapter.NewQueryAdapter(userQueryClient, g.User...),
		ReviewCommand: NewReviewCommandRepository(queries),
	}
}
