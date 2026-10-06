// Package product adapts the Product service gRPC API into the shared
// pkg/database sqlc row types. Query and command are exposed as separate
// interfaces so consumers can depend on the narrower one.
package product

import (
	"context"

	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	product_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
)

// QueryRepository reads products from the Product service over gRPC.
type QueryRepository interface {
	FindById(ctx context.Context, productID int) (*db.GetProductByIDRow, error)
	FindByID(ctx context.Context, productID int) (*db.GetProductByIDRow, error)
}

// CommandRepository mutates product stock via the Product service over gRPC.
type CommandRepository interface {
	UpdateProductCountStock(ctx context.Context, productID int, stock int) (*db.UpdateProductCountStockRow, error)
}

// Repository implements both QueryRepository and CommandRepository over the
// generated product query/command clients.
type Repository struct {
	query   pbproduct.ProductQueryServiceClient
	command pbproduct.ProductCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds a product adapter. Passing zero options leaves the guard nil, which
// makes DependencyGuard.Call a plain passthrough.
func New(query pbproduct.ProductQueryServiceClient, command pbproduct.ProductCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter builds a read-only adapter for consumers that never write.
func NewQueryAdapter(query pbproduct.ProductQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(query, nil, opts...)
}

// NewCommandAdapter builds a write-only adapter for consumers that never read.
func NewCommandAdapter(command pbproduct.ProductCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return New(nil, command, opts...)
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindById(ctx context.Context, productID int) (*db.GetProductByIDRow, error) {
	return r.findByID(ctx, productID)
}

func (r *Repository) FindByID(ctx context.Context, productID int) (*db.GetProductByIDRow, error) {
	return r.findByID(ctx, productID)
}

func (r *Repository) findByID(ctx context.Context, productID int) (*db.GetProductByIDRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbproduct.ApiResponseProduct, error) {
		return r.query.FindById(callCtx, &pbproduct.FindByIdProductRequest{Id: int32(productID)})
	})
	if err != nil {
		return nil, product_errors.ErrProductNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, product_errors.ErrProductNotFound
	}

	rating := float64(resp.Data.Rating)
	return &db.GetProductByIDRow{
		ProductID:    resp.Data.Id,
		MerchantID:   resp.Data.MerchantId,
		CategoryID:   resp.Data.CategoryId,
		Name:         resp.Data.Name,
		Description:  &resp.Data.Description,
		Price:        resp.Data.Price,
		CountInStock: resp.Data.CountInStock,
		Brand:        &resp.Data.Brand,
		Weight:       &resp.Data.Weight,
		Rating:       &rating,
		SlugProduct:  &resp.Data.SlugProduct,
		ImageProduct: &resp.Data.ImageProduct,
	}, nil
}

func (r *Repository) UpdateProductCountStock(ctx context.Context, productID int, stock int) (*db.UpdateProductCountStockRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbproduct.ApiResponseProduct, error) {
		return r.command.UpdateProductCountStock(callCtx, &pbproduct.UpdateProductCountStockRequest{
			ProductId: int32(productID),
			Stock:     int32(stock),
		})
	})
	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return &db.UpdateProductCountStockRow{
		ProductID:    resp.Data.Id,
		CountInStock: resp.Data.CountInStock,
	}, nil
}
