package repository

import (
	"context"

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

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps *Deps) *Repositories {
	mapperOrder := recordmapper.NewOrderRecordMapper()
	mapperOrderItem := recordmapper.NewOrderItemRecordMapper()
	mapperProduct := recordmapper.NewProductRecordMapper()
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperUser := recordmapper.NewUserRecordMapper()
	mapperShipping := recordmapper.NewShippingAddressRecordMapper()

	return &Repositories{
		MerchantQuery:    NewMerchantQueryRepository(deps.DB, deps.Ctx, mapperMerchant),
		ProductQuery:     NewProductQueryRepository(deps.DB, deps.Ctx, mapperProduct),
		ProductCommand:   NewProductCommandRepository(deps.DB, deps.Ctx, mapperProduct),
		OrderItemQuery:   NewOrderItemQueryRepository(deps.DB, deps.Ctx, mapperOrderItem),
		OrderItemCommand: NewOrderItemCommandRepository(deps.DB, deps.Ctx, mapperOrderItem),
		OrderQuery:       NewOrderQueryRepository(deps.DB, deps.Ctx, mapperOrder),
		OrderCommand:     NewOrderCommandRepository(deps.DB, deps.Ctx, mapperOrder),
		UserQuery:        NewUserQueryRepository(deps.DB, deps.Ctx, mapperUser),
		ShippingAddress:  NewShippingAddressCommandRepository(deps.DB, deps.Ctx, mapperShipping),
		OrderStats:       NewOrderStatsRepository(deps.DB, deps.Ctx, mapperOrder),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(
			deps.DB,
			deps.Ctx,
			mapperOrder,
		),
	}
}
