package repository

import (
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

func NewRepositories(DB *db.Queries) *Repositoris {
	mapperOrderItem := recordmapper.NewOrderItemRecordMapper()
	mapperOrder := recordmapper.NewOrderRecordMapper()
	mapperTransaction := recordmapper.NewTransactionRecordMapper()
	mapperShipping := recordmapper.NewShippingAddressRecordMapper()
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperUser := recordmapper.NewUserRecordMapper()

	return &Repositoris{
		TransactionCommandRepository:   NewTransactionCommandRepository(DB, mapperTransaction),
		TransactionQueryRepository:     NewTransactionQueryRepository(DB, mapperTransaction),
		OrderItemRepository:            NewOrderItemQueryRepository(DB, mapperOrderItem),
		OrderQueryRepository:           NewOrderQueryRepository(DB, mapperOrder),
		MerchantRepository:             NewMerchantQueryRepository(DB, mapperMerchant),
		ShippingAddressQueryRepository: NewShippingAddressQueryRepository(DB, mapperShipping),
		TransactionStatsRepository:     NewTransactionStatsRepository(DB, mapperTransaction),
		TransactionStatsByMerchant:     NewTransactionStatsByMerchantRepository(DB, mapperTransaction),
		UserQuery:                      NewUserQueryRepository(DB, mapperUser),
	}
}
