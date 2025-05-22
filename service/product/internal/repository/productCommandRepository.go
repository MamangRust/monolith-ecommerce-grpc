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
	ctx     context.Context
	mapping recordmapper.ProductRecordMapping
}

func NewProductCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.ProductRecordMapping) *productCommandRepository {
	return &productCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *productCommandRepository) CreateProduct(request *requests.CreateProductRequest) (*record.ProductRecord, error) {
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

	product, err := r.db.CreateProduct(r.ctx, req)

	if err != nil {
		return nil, product_errors.ErrCreateProduct
	}

	return r.mapping.ToProductRecord(product), nil
}

func (r *productCommandRepository) UpdateProduct(request *requests.UpdateProductRequest) (*record.ProductRecord, error) {
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

	res, err := r.db.UpdateProduct(r.ctx, req)

	if err != nil {
		return nil, product_errors.ErrUpdateProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) UpdateProductCountStock(product_id int, stock int) (*record.ProductRecord, error) {
	res, err := r.db.UpdateProductCountStock(r.ctx, db.UpdateProductCountStockParams{
		ProductID:    int32(product_id),
		CountInStock: int32(stock),
	})

	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) TrashedProduct(product_id int) (*record.ProductRecord, error) {
	res, err := r.db.TrashProduct(r.ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrTrashedProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) RestoreProduct(product_id int) (*record.ProductRecord, error) {
	res, err := r.db.RestoreProduct(r.ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrRestoreProduct
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productCommandRepository) DeleteProductPermanent(product_id int) (bool, error) {
	err := r.db.DeleteProductPermanently(r.ctx, int32(product_id))

	if err != nil {
		return false, product_errors.ErrDeleteProductPermanent
	}

	return true, nil
}

func (r *productCommandRepository) RestoreAllProducts() (bool, error) {
	err := r.db.RestoreAllProducts(r.ctx)

	if err != nil {
		return false, product_errors.ErrRestoreAllProducts
	}

	return true, nil
}

func (r *productCommandRepository) DeleteAllProductPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentProducts(r.ctx)

	if err != nil {
		return false, product_errors.ErrDeleteAllProductPermanent
	}

	return true, nil
}
