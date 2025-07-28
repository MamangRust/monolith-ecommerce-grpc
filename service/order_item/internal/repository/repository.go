package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	OrderItemQuery OrderItemQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapper := recordmapper.NewOrderItemRecordMapper()

	return &Repositories{
		OrderItemQuery: NewOrderItemQueryRepository(DB, mapper),
	}
}
