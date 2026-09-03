package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
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
	TransactionCommand   TransactionCommandRepository
	ShippingQuery        pb.ShippingQueryServiceClient
	StockReservation     StockReservationRepository
}

type Deps struct {
	DB               *db.Queries
	MerchantQuery    pb.MerchantQueryServiceClient
	ProductQuery     pb.ProductQueryServiceClient
	ProductCommand   pb.ProductCommandServiceClient
	OrderItemQuery   pb.OrderItemQueryServiceClient
	OrderItemCommand pb.OrderItemCommandServiceClient
	UserQuery        pb.UserQueryServiceClient
	ShippingCommand  pb.ShippingCommandServiceClient
	TransactionCommand pb.TransactionCommandServiceClient
	ShippingQuery      pb.ShippingQueryServiceClient
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		MerchantQuery:    NewMerchantQueryRepository(deps.MerchantQuery),
		ProductQuery:     NewProductQueryRepository(deps.ProductQuery),
		ProductCommand:   NewProductCommandRepository(deps.ProductCommand),
		OrderItemQuery:   NewOrderItemQueryRepository(deps.OrderItemQuery, deps.OrderItemCommand),
		OrderItemCommand: NewOrderItemCommandRepository(deps.OrderItemCommand),
		OrderQuery:       NewOrderQueryRepository(deps.DB),
		OrderCommand:     NewOrderCommandRepository(deps.DB),
		UserQuery:        NewUserQueryRepository(deps.UserQuery),
		ShippingAddress:  NewShippingAddressCommandRepository(deps.ShippingCommand),
		TransactionCommand: NewTransactionCommandRepository(deps.TransactionCommand),
		ShippingQuery:      deps.ShippingQuery,
		StockReservation:   NewStockReservationRepository(deps.DB),
		OrderStats:       NewOrderStatsRepository(deps.DB),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(
			deps.DB,
		),
	}
}
