package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	OrderItemQuery OrderItemQueryRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps *Deps) *Repositories {
	mapper := recordmapper.NewOrderItemRecordMapper()

	return &Repositories{
		OrderItemQuery: NewOrderItemQueryRepository(deps.DB, deps.Ctx, mapper),
	}
}
