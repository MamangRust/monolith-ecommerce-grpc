package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	ShippingAddressQuery   ShippingAddressQueryRepository
	ShippingAddressCommand ShippingAddressCommandRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps Deps) *Repositories {
	mapper := recordmapper.NewShippingAddressRecordMapper()

	return &Repositories{
		ShippingAddressQuery:   NewShippingAddressQueryRepository(deps.DB, deps.Ctx, mapper),
		ShippingAddressCommand: NewShippingAddressCommandRepository(deps.DB, deps.Ctx, mapper),
	}
}
