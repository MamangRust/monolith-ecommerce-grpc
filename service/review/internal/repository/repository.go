package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	ProductQuery  ProductQueryRepository
	ReviewQuery   ReviewQueryRepository
	UserQuery     UserQueryRepository
	ReviewCommand ReviewCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapperProduct := recordmapper.NewProductRecordMapper()
	mapperReview := recordmapper.NewReviewRecordMapper()
	mapperUser := recordmapper.NewUserRecordMapper()

	return &Repositories{
		ProductQuery:  NewProductQueryRepository(DB, mapperProduct),
		ReviewQuery:   NewReviewQueryRepository(DB, mapperReview),
		UserQuery:     NewUserQueryRepository(DB, mapperUser),
		ReviewCommand: NewReviewCommandRepository(DB, mapperReview),
	}
}
