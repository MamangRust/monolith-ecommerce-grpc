package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositoris struct {
	TransactionCommandRepository   TransactionCommandRepository
	TransactionQueryRepository     TransactionQueryRepository
	OrderItemRepository            OrderItemRepository
	OrderQueryRepository           OrderQueryRepository
	MerchantRepository             MerchantQueryRepository
	ShippingAddressQueryRepository ShippingAddressQueryRepository
	TransactionStatsRepository     TransactionStatsRepository
	TransactionStatsByMerchant     TransactonStatsByMerchantRepository
	UserQuery                      UserQueryRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps Deps) *Repositoris {
	mapperOrderItem := recordmapper.NewOrderItemRecordMapper()
	mapperOrder := recordmapper.NewOrderRecordMapper()
	mapperTransaction := recordmapper.NewTransactionRecordMapper()
	mapperShipping := recordmapper.NewShippingAddressRecordMapper()
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperUser := recordmapper.NewUserRecordMapper()

	return &Repositoris{
		TransactionCommandRepository:   NewTransactionCommandRepository(deps.DB, deps.Ctx, mapperTransaction),
		TransactionQueryRepository:     NewTransactionQueryRepository(deps.DB, deps.Ctx, mapperTransaction),
		OrderItemRepository:            NewOrderItemQueryRepository(deps.DB, deps.Ctx, mapperOrderItem),
		OrderQueryRepository:           NewOrderQueryRepository(deps.DB, deps.Ctx, mapperOrder),
		MerchantRepository:             NewMerchantQueryRepository(deps.DB, deps.Ctx, mapperMerchant),
		ShippingAddressQueryRepository: NewShippingAddressQueryRepository(deps.DB, deps.Ctx, mapperShipping),
		TransactionStatsRepository:     NewTransactionStatsRepository(deps.DB, deps.Ctx, mapperTransaction),
		TransactionStatsByMerchant:     NewTransactionStatsByMerchantRepository(deps.DB, deps.Ctx, mapperTransaction),
		UserQuery:                      NewUserQueryRepository(deps.DB, deps.Ctx, mapperUser),
	}
}
