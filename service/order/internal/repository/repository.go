package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
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
	mapperOrder := recordmapper.NewOrderRecordMapper()
	mapperOrderItem := recordmapper.NewOrderItemRecordMapper()
	mapperProduct := recordmapper.NewProductRecordMapper()
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperUser := recordmapper.NewUserRecordMapper()
	mapperShipping := recordmapper.NewShippingAddressRecordMapper()

	return &Repositories{
		MerchantQuery:    NewMerchantQueryRepository(DB, mapperMerchant),
		ProductQuery:     NewProductQueryRepository(DB, mapperProduct),
		ProductCommand:   NewProductCommandRepository(DB, mapperProduct),
		OrderItemQuery:   NewOrderItemQueryRepository(DB, mapperOrderItem),
		OrderItemCommand: NewOrderItemCommandRepository(DB, mapperOrderItem),
		OrderQuery:       NewOrderQueryRepository(DB, mapperOrder),
		OrderCommand:     NewOrderCommandRepository(DB, mapperOrder),
		UserQuery:        NewUserQueryRepository(DB, mapperUser),
		ShippingAddress:  NewShippingAddressCommandRepository(DB, mapperShipping),
		OrderStats:       NewOrderStatsRepository(DB, mapperOrder),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(
			DB,
			mapperOrder,
		),
	}
}
