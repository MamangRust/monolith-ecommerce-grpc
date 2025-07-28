package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	ShippingAddressQuery   ShippingAddressQueryRepository
	ShippingAddressCommand ShippingAddressCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapper := recordmapper.NewShippingAddressRecordMapper()

	return &Repositories{
		ShippingAddressQuery:   NewShippingAddressQueryRepository(DB, mapper),
		ShippingAddressCommand: NewShippingAddressCommandRepository(DB, mapper),
	}
}
