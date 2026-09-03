package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/convert"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
)

type productCommandRepository struct {
	db *db.Queries
}

func NewProductCommandRepository(db *db.Queries) *productCommandRepository {
	return &productCommandRepository{
		db: db,
	}
}

func (r *productCommandRepository) Create(ctx context.Context, request *requests.CreateProductRequest) (*db.CreateProductRow, error) {
	req := db.CreateProductParams{
		MerchantID:   int32(request.MerchantID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  convert.StringPtr(request.Description),
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        convert.StringPtr(request.Brand),
		Weight:       int32Ptr(request.Weight),
		SlugProduct:  request.SlugProduct,
		ImageProduct: convert.StringPtr(request.ImageProduct),
	}

	product, err := r.db.CreateProduct(ctx, req)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product_errors.ErrProductNotFound
		}
		return nil, product_errors.ErrCreateProduct.WithInternal(err)
	}

	return product, nil
}

func (r *productCommandRepository) Update(ctx context.Context, request *requests.UpdateProductRequest) (*db.UpdateProductRow, error) {
	req := db.UpdateProductParams{
		ProductID:    int32(*request.ProductID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  convert.StringPtr(request.Description),
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        convert.StringPtr(request.Brand),
		Weight:       int32Ptr(request.Weight),
		SlugProduct:  request.SlugProduct,
		ImageProduct: convert.StringPtr(request.ImageProduct),
	}

	res, err := r.db.UpdateProduct(ctx, req)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product_errors.ErrProductNotFound
		}
		return nil, product_errors.ErrUpdateProduct.WithInternal(err)
	}

	return res, nil
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.UpdateProductCountStockRow, error) {
	res, err := r.db.UpdateProductCountStock(ctx, db.UpdateProductCountStockParams{
		ProductID:    int32(product_id),
		CountInStock: int32(stock),
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product_errors.ErrProductNotFound
		}
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return res, nil
}

func (r *productCommandRepository) AdjustProductStock(ctx context.Context, product_id int, delta int, operationID string) (*db.AdjustProductStockRow, error) {
	res, err := r.db.AdjustProductStock(ctx, db.AdjustProductStockParams{
		OperationID: operationID,
		ProductID:   int32(product_id),
		Delta:       int32(delta),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product_errors.ErrProductNotFound
		}
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}

	return res, nil
}

func (r *productCommandRepository) Trash(ctx context.Context, product_id int) (*db.TrashProductRow, error) {
	res, err := r.db.TrashProduct(ctx, int32(product_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product_errors.ErrProductNotFound
		}
		return nil, product_errors.ErrTrashedProduct.WithInternal(err)
	}

	return res, nil
}

func (r *productCommandRepository) Restore(ctx context.Context, product_id int) (*db.RestoreProductRow, error) {
	res, err := r.db.RestoreProduct(ctx, int32(product_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product_errors.ErrProductNotFound
		}
		return nil, product_errors.ErrRestoreProduct.WithInternal(err)
	}

	return res, nil
}

func (r *productCommandRepository) DeletePermanent(ctx context.Context, product_id int) (bool, error) {
	err := r.db.DeleteProductPermanently(ctx, int32(product_id))

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return false, shared_errors.NewConflictError("cannot permanently delete product while related records exist").WithInternal(err)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return false, product_errors.ErrProductNotFound
		}
		return false, product_errors.ErrDeleteProductPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *productCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllProducts(ctx)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, product_errors.ErrProductNotFound
		}
		return false, product_errors.ErrRestoreAllProducts.WithInternal(err)
	}

	return true, nil
}

func (r *productCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentProducts(ctx)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return false, shared_errors.NewConflictError("cannot permanently delete products while related records exist").WithInternal(err)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return false, product_errors.ErrProductNotFound
		}
		return false, product_errors.ErrDeleteAllProducts.WithInternal(err)
	}

	return true, nil
}
