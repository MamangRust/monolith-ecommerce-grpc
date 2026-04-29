package repository

import (
	"context"

	"database/sql"
 
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
)


type productQueryRepository struct {
	db *db.Queries
}

func NewProductQueryRepository(db *db.Queries) *productQueryRepository {
	return &productQueryRepository{
		db: db,
	}
}

func (r *productQueryRepository) FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProducts(ctx, reqDb)
	if err != nil {
		return nil, product_errors.ErrFindAllProducts.WithInternal(err)
	}


	return res, nil
}

func (r *productQueryRepository) FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsActive(ctx, reqDb)
	if err != nil {
		return nil, product_errors.ErrFindActiveProducts.WithInternal(err)
	}


	return res, nil
}

func (r *productQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsTrashed(ctx, reqDb)
	if err != nil {
		return nil, product_errors.ErrFindTrashedProducts.WithInternal(err)
	}


	return res, nil
}

func (r *productQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    stringPtr(req.Search),
		Column3:    int32(req.CategoryID),
		Column4:    int32(IntPtrToInt(req.MinPrice)),
		Column5:    int32(IntPtrToInt(req.MaxPrice)),
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetProductsByMerchant(ctx, reqDb)
	if err != nil {
		return nil, product_errors.ErrFindProductsByMerchant.WithInternal(err)
	}


	return res, nil
}

func (r *productQueryRepository) FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsByCategoryNameParams{
		Name:    req.CategoryName,
		Column2: req.Search,
		Column3: int32(IntPtrToInt(req.MinPrice)),
		Column4: int32(IntPtrToInt(req.MaxPrice)),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsByCategoryName(ctx, reqDb)
	if err != nil {
		return nil, product_errors.ErrFindProductsByCategory.WithInternal(err)
	}


	return res, nil
}

func (r *productQueryRepository) FindByID(ctx context.Context, product_id int) (*db.GetProductByIDRow, error) {
	res, err := r.db.GetProductByID(ctx, int32(product_id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, product_errors.ErrProductNotFound.WithInternal(err)
		}
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}


	return res, nil
}

func (r *productQueryRepository) FindByIDTrashed(ctx context.Context, product_id int) (*db.Product, error) {
	res, err := r.db.GetProductByIdTrashed(ctx, int32(product_id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, product_errors.ErrProductNotFound.WithInternal(err)
		}
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}


	return res, nil
}
