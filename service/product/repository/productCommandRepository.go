package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
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
		Description:  stringPtr(request.Description),
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        stringPtr(request.Brand),
		Weight:       int32Ptr(request.Weight),
		SlugProduct:  request.SlugProduct,
		ImageProduct: stringPtr(request.ImageProduct),
	}

	product, err := r.db.CreateProduct(ctx, req)

	if err != nil {
		return nil, product_errors.ErrCreateProduct.WithInternal(err)
	}


	return product, nil
}

func (r *productCommandRepository) Update(ctx context.Context, request *requests.UpdateProductRequest) (*db.UpdateProductRow, error) {
	req := db.UpdateProductParams{
		ProductID:    int32(*request.ProductID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  stringPtr(request.Description),
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        stringPtr(request.Brand),
		Weight:       int32Ptr(request.Weight),
		SlugProduct:  request.SlugProduct,
		ImageProduct: stringPtr(request.ImageProduct),
	}

	res, err := r.db.UpdateProduct(ctx, req)

	if err != nil {
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
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}


	return res, nil
}

func (r *productCommandRepository) Trash(ctx context.Context, product_id int) (*db.TrashProductRow, error) {
	res, err := r.db.TrashProduct(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrTrashedProduct.WithInternal(err)
	}


	return res, nil
}

func (r *productCommandRepository) Restore(ctx context.Context, product_id int) (*db.RestoreProductRow, error) {
	res, err := r.db.RestoreProduct(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrRestoreProduct.WithInternal(err)
	}


	return res, nil
}

func (r *productCommandRepository) DeletePermanent(ctx context.Context, product_id int) (bool, error) {
	err := r.db.DeleteProductPermanently(ctx, int32(product_id))

	if err != nil {
		return false, product_errors.ErrDeleteProductPermanent.WithInternal(err)
	}


	return true, nil
}

func (r *productCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllProducts(ctx)

	if err != nil {
		return false, product_errors.ErrRestoreAllProducts.WithInternal(err)
	}


	return true, nil
}

func (r *productCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentProducts(ctx)

	if err != nil {
		return false, product_errors.ErrDeleteAllProducts.WithInternal(err)
	}


	return true, nil
}

