package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type Repositories struct {
	CartQuery    CartQueryRepository
	CartCommand  CartCommandRepository
	UserQuery    UserQueryRepository
	ProductQuery ProductQueryRepository
}

func NewRepositories(DB *db.Queries,
	userQuery pb.UserQueryServiceClient,
	productQuery pb.ProductQueryServiceClient,
) *Repositories {
	return &Repositories{
		CartQuery:    NewCartQueryRepository(DB),
		CartCommand:  NewCartCommandRepository(DB),
		UserQuery:    NewUserQueryRepository(userQuery),
		ProductQuery: NewProductQueryRepository(productQuery),
	}
}
