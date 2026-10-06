package repository

import (
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pborder "github.com/MamangRust/monolith-ecommerce-pb/order"
	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
	pbshipping_address "github.com/MamangRust/monolith-ecommerce-pb/shipping_address"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/order_item"
	shippingaddressadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/shipping_address"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
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
	Outbox             OutboxRepository
}

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	User      []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Order     []adapter.GuardOption
	OrderItem []adapter.GuardOption
	Shipping  []adapter.GuardOption
}

type Deps struct {
	Db *db.Queries

	UserQueryClient      pbuser.UserQueryServiceClient
	MerchantQueryClient  pbmerchant.MerchantQueryServiceClient
	OrderQueryClient     pborder.OrderQueryServiceClient
	OrderItemQueryClient pborder_item.OrderItemQueryServiceClient
	ShippingQueryClient  pbshipping_address.ShippingQueryServiceClient

	Guards GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	g := deps.Guards

	return &Repositories{
		TransactionCommand: NewTransactionCommandRepository(deps.Db),
		TransactionQuery:   NewTransactionQueryRepository(deps.Db),
		OrderItem:          orderitemadapter.New(deps.OrderItemQueryClient, nil, g.OrderItem...),
		OrderQuery:         orderadapter.New(deps.OrderQueryClient, g.Order...),
		MerchantQuery:      merchantadapter.New(deps.MerchantQueryClient, g.Merchant...),
		ShippingAddress:    shippingaddressadapter.New(deps.ShippingQueryClient, nil, g.Shipping...),
		TransactionStats:   NewTransactionStatsRepository(deps.Db),
		StatsByMerchant:    NewTransactionStatsByMerchantRepository(deps.Db),
		UserQuery:          useradapter.New(deps.UserQueryClient, nil, g.User...),
		Outbox:             NewOutboxRepository(deps.Db),
	}
}
