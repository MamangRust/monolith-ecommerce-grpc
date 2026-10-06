package repository

import (
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
	pbshipping_address "github.com/MamangRust/monolith-ecommerce-pb/shipping_address"
	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	orderitemadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/order_item"
	productadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/product"
	shippingaddressadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/shipping_address"
	transactionadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/transaction"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
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
	TransactionCommand   TransactionCommandRepository
	ShippingQuery        shippingaddressadapter.QueryRepository
}

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	Merchant    []adapter.GuardOption
	Product     []adapter.GuardOption
	OrderItem   []adapter.GuardOption
	User        []adapter.GuardOption
	Shipping    []adapter.GuardOption
	Transaction []adapter.GuardOption
}

type Deps struct {
	Db *db.Queries

	MerchantQueryClient      pbmerchant.MerchantQueryServiceClient
	ProductQueryClient       pbproduct.ProductQueryServiceClient
	ProductCommandClient     pbproduct.ProductCommandServiceClient
	OrderItemQueryClient     pborder_item.OrderItemQueryServiceClient
	OrderItemCommandClient   pborder_item.OrderItemCommandServiceClient
	UserQueryClient          pbuser.UserQueryServiceClient
	ShippingCommandClient    pbshipping_address.ShippingCommandServiceClient
	ShippingQueryClient      pbshipping_address.ShippingQueryServiceClient
	TransactionCommandClient pbtransaction.TransactionCommandServiceClient

	Guards GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	g := deps.Guards

	productAdapter := productadapter.New(deps.ProductQueryClient, deps.ProductCommandClient, g.Product...)
	orderItemAdapter := orderitemadapter.New(deps.OrderItemQueryClient, deps.OrderItemCommandClient, g.OrderItem...)
	shippingAdapter := shippingaddressadapter.New(deps.ShippingQueryClient, deps.ShippingCommandClient, g.Shipping...)

	return &Repositories{
		MerchantQuery:        merchantadapter.New(deps.MerchantQueryClient, g.Merchant...),
		ProductQuery:         productAdapter,
		ProductCommand:       productAdapter,
		OrderItemQuery:       orderItemAdapter,
		OrderItemCommand:     orderItemAdapter,
		OrderQuery:           NewOrderQueryRepository(deps.Db),
		OrderCommand:         NewOrderCommandRepository(deps.Db),
		UserQuery:            useradapter.New(deps.UserQueryClient, nil, g.User...),
		ShippingAddress:      shippingAdapter,
		TransactionCommand:   transactionadapter.New(deps.TransactionCommandClient, g.Transaction...),
		ShippingQuery:        shippingAdapter,
		OrderStats:           NewOrderStatsRepository(deps.Db),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(deps.Db),
	}
}
