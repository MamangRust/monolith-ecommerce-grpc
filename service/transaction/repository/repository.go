package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
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

type Deps struct {
	DB              *db.Queries
	UserQuery       pb.UserQueryServiceClient
	MerchantQuery   pb.MerchantQueryServiceClient
	OrderQuery      pb.OrderQueryServiceClient
	OrderItemQuery  pb.OrderItemQueryServiceClient
	ShippingQuery   pb.ShippingQueryServiceClient
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		TransactionCommand: NewTransactionCommandRepository(deps.DB),
		TransactionQuery:   NewTransactionQueryRepository(deps.DB),
		OrderItem:          NewOrderItemRepository(deps.OrderItemQuery),
		OrderQuery:         NewOrderQueryRepository(deps.OrderQuery),
		MerchantQuery:      NewMerchantQueryRepository(deps.MerchantQuery),
		ShippingAddress:    NewShippingAddressQueryRepository(deps.ShippingQuery),
		TransactionStats:   NewTransactionStatsRepository(deps.DB),
		StatsByMerchant:    NewTransactionStatsByMerchantRepository(deps.DB),
		UserQuery:          NewUserQueryRepository(deps.UserQuery),
	}
}
