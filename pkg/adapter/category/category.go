// Package category adapts the Category service gRPC API into the shared
// pkg/database sqlc row types.
package category

import (
	"context"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
)

// QueryRepository is the contract consumers depend on for category reads.
type QueryRepository interface {
	FindByID(ctx context.Context, categoryID int) (*db.GetCategoryByIDRow, error)
}

type repository struct {
	client pbcategory.CategoryQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a category adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough.
func New(client pbcategory.CategoryQueryServiceClient, opts ...adapter.GuardOption) *repository {
	r := &repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter returns the adapter typed as the read-only contract.
func NewQueryAdapter(client pbcategory.CategoryQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(client, opts...)
}

func (r *repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *repository) FindByID(ctx context.Context, categoryID int) (*db.GetCategoryByIDRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbcategory.ApiResponseCategory, error) {
		return r.client.FindById(callCtx, &pbcategory.FindByIdCategoryRequest{Id: int32(categoryID)})
	})
	if err != nil {
		return nil, category_errors.ErrFindCategoryById.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, category_errors.ErrFindCategoryById
	}

	return &db.GetCategoryByIDRow{
		CategoryID:  resp.Data.Id,
		Name:        resp.Data.Name,
		Description: &resp.Data.Description,
	}, nil
}
