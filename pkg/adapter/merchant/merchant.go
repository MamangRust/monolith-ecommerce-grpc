// Package merchant adapts the Merchant service gRPC API into the shared
// pkg/database sqlc row types used by the merchant-domain consumers.
package merchant

import (
	"context"

	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
)

// QueryRepository is the contract consumers depend on for merchant reads.
type QueryRepository interface {
	FindByID(ctx context.Context, merchantID int) (*db.GetMerchantByIDRow, error)
}

type repository struct {
	client pbmerchant.MerchantQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a merchant adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough.
func New(client pbmerchant.MerchantQueryServiceClient, opts ...adapter.GuardOption) *repository {
	r := &repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter returns the adapter typed as the read-only contract.
func NewQueryAdapter(client pbmerchant.MerchantQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(client, opts...)
}

func (r *repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *repository) FindByID(ctx context.Context, merchantID int) (*db.GetMerchantByIDRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbmerchant.ApiResponseMerchant, error) {
		return r.client.FindById(callCtx, &pbmerchant.FindByIdMerchantRequest{Id: int32(merchantID)})
	})
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	return &db.GetMerchantByIDRow{
		MerchantID:   resp.Data.Id,
		UserID:       resp.Data.UserId,
		Name:         resp.Data.Name,
		Description:  &resp.Data.Description,
		Address:      &resp.Data.Address,
		ContactEmail: &resp.Data.ContactEmail,
		ContactPhone: &resp.Data.ContactPhone,
		Status:       resp.Data.Status,
	}, nil
}
