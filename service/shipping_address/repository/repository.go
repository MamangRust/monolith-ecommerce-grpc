package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

type Repositories struct {
	ShippingAddressQuery   ShippingAddressQueryRepository
	ShippingAddressCommand ShippingAddressCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		ShippingAddressQuery:   NewShippingAddressQueryRepository(DB),
		ShippingAddressCommand: NewShippingAddressCommandRepository(DB),
	}
}
