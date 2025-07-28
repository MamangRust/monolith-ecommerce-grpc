package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	ReviewDetailQuery   ReviewDetailQueryRepository
	ReviewDetailCommand ReviewDetailCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapper := recordmapper.NewReviewDetailRecordMapper()

	return &Repositories{
		ReviewDetailQuery:   NewReviewDetailQueryRepository(DB, mapper),
		ReviewDetailCommand: NewReviewDetailCommandRepository(DB, mapper),
	}
}
