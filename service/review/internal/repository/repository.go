package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	ProductQuery  ProductQueryRepository
	ReviewQuery   ReviewQueryRepository
	UserQuery     UserQueryRepository
	ReviewCommand ReviewCommandRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps Deps) *Repositories {
	mapperProduct := recordmapper.NewProductRecordMapper()
	mapperReview := recordmapper.NewReviewRecordMapper()
	mapperUser := recordmapper.NewUserRecordMapper()

	return &Repositories{
		ProductQuery:  NewProductQueryRepository(deps.DB, deps.Ctx, mapperProduct),
		ReviewQuery:   NewReviewQueryRepository(deps.DB, deps.Ctx, mapperReview),
		UserQuery:     NewUserQueryRepository(deps.DB, deps.Ctx, mapperUser),
		ReviewCommand: NewReviewCommandRepository(deps.DB, deps.Ctx, mapperReview),
	}
}
