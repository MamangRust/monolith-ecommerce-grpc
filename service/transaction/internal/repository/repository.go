package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

type Repositories struct {
	TransactionCommand TransactionCommandRepository
	TransactionQuery   TransactionQueryRepository
	OrderItem          OrderItemRepository
	OrderQuery         OrderQueryRepository
	MerchantQuery      MerchantQueryRepository
	ShippingAddress    ShippingAddressQueryRepository
	TransactionStats   TransactionStatsRepository
	StatsByMerchant    TransactionStatsByMerchantRepository
	UserQuery          UserQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		TransactionCommand: NewTransactionCommandRepository(DB),
		TransactionQuery:   NewTransactionQueryRepository(DB),
		OrderItem:          NewOrderItemQueryRepository(DB),
		OrderQuery:         NewOrderQueryRepository(DB),
		MerchantQuery:      NewMerchantQueryRepository(DB),
		ShippingAddress:    NewShippingAddressQueryRepository(DB),
		TransactionStats:   NewTransactionStatsRepository(DB),
		StatsByMerchant:    NewTransactionStatsByMerchantRepository(DB),
		UserQuery:          NewUserQueryRepository(DB),
	}
}
