package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

type Repositories struct {
	MerchantQuery        MerchantQueryRepository
	ProductQuery         ProductQueryRepository
	ProductCommand       ProductCommandRepository
	OrderItemQuery       OrderItemQueryRepository
	OrderItemCommand     OrderItemCommandRepository
	OrderQuery           OrderQueryRepository
	OrderCommand         OrderCommandRepository
	OrderStats           OrderStatsRepository
	OrderStatsByMerchant OrderStatsByMerchantRepository
	UserQuery            UserQueryRepository
	ShippingAddress      ShippingAddressCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		MerchantQuery:    NewMerchantQueryRepository(DB),
		ProductQuery:     NewProductQueryRepository(DB),
		ProductCommand:   NewProductCommandRepository(DB),
		OrderItemQuery:   NewOrderItemQueryRepository(DB),
		OrderItemCommand: NewOrderItemCommandRepository(DB),
		OrderQuery:       NewOrderQueryRepository(DB),
		OrderCommand:     NewOrderCommandRepository(DB),
		UserQuery:        NewUserQueryRepository(DB),
		ShippingAddress:  NewShippingAddressCommandRepository(DB),
		OrderStats:       NewOrderStatsRepository(DB),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(
			DB,
		),
	}
}
