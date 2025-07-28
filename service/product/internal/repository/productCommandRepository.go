package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type productCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.ProductRecordMapping
}

func NewProductCommandRepository(db *db.Queries, mapping recordmapper.ProductRecordMapping) *productCommandRepository {
	return &productCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *productCommandRepository) CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*record.ProductRecord, error) {
	req := db.CreateProductParams{
		MerchantID:   int32(request.MerchantID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        sql.NullString{String: request.Brand, Valid: true},
		Weight:       sql.NullInt32{Int32: int32(request.Weight), Valid: true},
		SlugProduct: sql.NullString{
			String: *request.SlugProduct,
			Valid:  true,
		},
		ImageProduct: sql.NullString{String: request.ImageProduct, Valid: true},
	}

	product, err := r.db.CreateProduct(ctx, req)

	if err != nil {
		return nil, product_errors.ErrCreateProduct
	}

	return r.mapping.ToProductRecord(product), nil
}

func (r *productCommandRepository) UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*record.ProductRecord, error) {
	req := db.UpdateProductParams{
		ProductID:    int32(*request.ProductID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        sql.NullString{String: request.Brand, Valid: true},
		Weight:       sql.NullInt32{Int32: int32(request.Weight), Valid: true},
		ImageProduct: sql.NullString{String: request.ImageProduct, Valid: true},
	}

	res, err := r.db.UpdateProduct(ctx, req)

	if err != nil {
		return nil, product_errors.ErrUpdateProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*record.ProductRecord, error) {
	res, err := r.db.UpdateProductCountStock(ctx, db.UpdateProductCountStockParams{
		ProductID:    int32(product_id),
		CountInStock: int32(stock),
	})

	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) TrashedProduct(ctx context.Context, product_id int) (*record.ProductRecord, error) {
	res, err := r.db.TrashProduct(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrTrashedProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) RestoreProduct(ctx context.Context, product_id int) (*record.ProductRecord, error) {
	res, err := r.db.RestoreProduct(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrRestoreProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) DeleteProductPermanent(ctx context.Context, product_id int) (bool, error) {
	err := r.db.DeleteProductPermanently(ctx, int32(product_id))

	if err != nil {
		return false, product_errors.ErrDeleteProductPermanent
	}

	return true, nil
}

func (r *productCommandRepository) RestoreAllProducts(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllProducts(ctx)

	if err != nil {
		return false, product_errors.ErrRestoreAllProducts
	}

	return true, nil
}

func (r *productCommandRepository) DeleteAllProductPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentProducts(ctx)

	if err != nil {
		return false, product_errors.ErrDeleteAllProductPermanent
	}

	return true, nil
}
