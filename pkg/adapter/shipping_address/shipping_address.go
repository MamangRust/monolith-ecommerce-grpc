// Package shipping_address adapts the Shipping Address service gRPC API into
// the shared pkg/database sqlc row types.
package shipping_address

import (
	"context"

	pbshipping "github.com/MamangRust/monolith-ecommerce-pb/shipping_address"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// QueryRepository is the contract consumers depend on for shipping-address reads.
type QueryRepository interface {
	FindByID(ctx context.Context, shippingID int) (*db.GetShippingAddressByOrderIDRow, error)
	FindByOrder(ctx context.Context, orderID int) (*db.GetShippingAddressByOrderIDRow, error)
}

// CommandRepository is the contract consumers depend on for shipping-address writes.
type CommandRepository interface {
	Create(ctx context.Context, req *requests.CreateShippingAddressRequest) (*db.CreateShippingAddressRow, error)
	Update(ctx context.Context, req *requests.UpdateShippingAddressRequest) (*db.UpdateShippingAddressRow, error)
	DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// Repository implements both QueryRepository and CommandRepository on top of the
// generated shipping query/command clients.
type Repository struct {
	query   pbshipping.ShippingQueryServiceClient
	command pbshipping.ShippingCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds a shipping-address adapter. Passing zero options leaves the guard
// nil, which makes DependencyGuard.Call a plain passthrough.
func New(query pbshipping.ShippingQueryServiceClient, command pbshipping.ShippingCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter builds a read-only adapter for consumers that never write.
func NewQueryAdapter(query pbshipping.ShippingQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(query, nil, opts...)
}

// NewCommandAdapter builds a write-only adapter for consumers that never read.
func NewCommandAdapter(command pbshipping.ShippingCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return New(nil, command, opts...)
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindByID(ctx context.Context, shippingID int) (*db.GetShippingAddressByOrderIDRow, error) {
	return r.find(ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return r.query.FindById(callCtx, &pbshipping.FindByIdShippingRequest{Id: int32(shippingID)})
	})
}

func (r *Repository) FindByOrder(ctx context.Context, orderID int) (*db.GetShippingAddressByOrderIDRow, error) {
	return r.find(ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return r.query.FindByOrder(callCtx, &pbshipping.FindByIdShippingRequest{Id: int32(orderID)})
	})
}

func (r *Repository) find(ctx context.Context, call func(context.Context) (*pbshipping.ApiResponseShipping, error)) (*db.GetShippingAddressByOrderIDRow, error) {
	resp, err := adapter.Call(r.guard, ctx, call)
	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder
	}

	return &db.GetShippingAddressByOrderIDRow{
		ShippingAddressID: resp.Data.Id,
		OrderID:           resp.Data.OrderId,
		Alamat:            resp.Data.Alamat,
		Provinsi:          resp.Data.Provinsi,
		Negara:            resp.Data.Negara,
		Kota:              resp.Data.Kota,
		ShippingMethod:    resp.Data.ShippingMethod,
		ShippingCost:      float64(resp.Data.ShippingCost),
	}, nil
}

func (r *Repository) Create(ctx context.Context, req *requests.CreateShippingAddressRequest) (*db.CreateShippingAddressRow, error) {
	var orderID int32
	if req.OrderID != nil {
		orderID = int32(*req.OrderID)
	}

	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return r.command.CreateShipping(callCtx, &pbshipping.CreateShippingAddressRequest{
			OrderId:        orderID,
			Alamat:         req.Alamat,
			Provinsi:       req.Provinsi,
			Kota:           req.Kota,
			Negara:         req.Negara,
			Courier:        req.Courier,
			ShippingMethod: req.ShippingMethod,
			ShippingCost:   int32(req.ShippingCost),
		})
	})
	if err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress
	}

	return &db.CreateShippingAddressRow{
		ShippingAddressID: resp.Data.Id,
		OrderID:           resp.Data.OrderId,
		Alamat:            resp.Data.Alamat,
		Provinsi:          resp.Data.Provinsi,
		Negara:            resp.Data.Negara,
		Kota:              resp.Data.Kota,
		ShippingMethod:    resp.Data.ShippingMethod,
		ShippingCost:      float64(resp.Data.ShippingCost),
	}, nil
}

func (r *Repository) Update(ctx context.Context, req *requests.UpdateShippingAddressRequest) (*db.UpdateShippingAddressRow, error) {
	var shippingID int32
	if req.ShippingID != nil {
		shippingID = int32(*req.ShippingID)
	}

	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return r.command.UpdateShipping(callCtx, &pbshipping.UpdateShippingAddressRequest{
			ShippingId:     shippingID,
			Alamat:         req.Alamat,
			Provinsi:       req.Provinsi,
			Kota:           req.Kota,
			Negara:         req.Negara,
			Courier:        req.Courier,
			ShippingMethod: req.ShippingMethod,
			ShippingCost:   int32(req.ShippingCost),
		})
	})
	if err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress
	}

	return &db.UpdateShippingAddressRow{
		ShippingAddressID: resp.Data.Id,
		OrderID:           resp.Data.OrderId,
		Alamat:            resp.Data.Alamat,
		Provinsi:          resp.Data.Provinsi,
		Negara:            resp.Data.Negara,
		Kota:              resp.Data.Kota,
		ShippingMethod:    resp.Data.ShippingMethod,
		ShippingCost:      float64(resp.Data.ShippingCost),
	}, nil
}

func (r *Repository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	_, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShippingDelete, error) {
		return r.command.DeleteShippingByOrderPermanent(callCtx, &pbshipping.FindByIdShippingRequest{
			Id: int32(orderID),
		})
	})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent.WithInternal(err)
	}
	return true, nil
}

func (r *Repository) DeleteAll(ctx context.Context) (bool, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShippingAll, error) {
		return r.command.DeleteAllShippingPermanent(callCtx, &emptypb.Empty{})
	})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteAllPermanentShippingAddress.WithInternal(err)
	}
	return resp != nil && resp.Status == "success", nil
}
