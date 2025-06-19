package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	ReviewDetailQuery   ReviewDetailQueryRepository
	ReviewDetailCommand ReviewDetailCommandRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps *Deps) *Repositories {
	mapper := recordmapper.NewReviewDetailRecordMapper()

	return &Repositories{
		ReviewDetailQuery:   NewReviewDetailQueryRepository(deps.DB, deps.Ctx, mapper),
		ReviewDetailCommand: NewReviewDetailCommandRepository(deps.DB, deps.Ctx, mapper),
	}
}
