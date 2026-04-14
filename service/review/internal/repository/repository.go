package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

type Repositories struct {
	ProductQuery  ProductQueryRepository
	ReviewQuery   ReviewQueryRepository
	UserQuery     UserQueryRepository
	ReviewCommand ReviewCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		ProductQuery:  NewProductQueryRepository(DB),
		ReviewQuery:   NewReviewQueryRepository(DB),
		UserQuery:     NewUserQueryRepository(DB),
		ReviewCommand: NewReviewCommandRepository(DB),
	}
}
