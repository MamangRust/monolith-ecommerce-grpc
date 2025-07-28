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

type productQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.ProductRecordMapping
}

func NewProductQueryRepository(db *db.Queries, mapping recordmapper.ProductRecordMapping) *productQueryRepository {
	return &productQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *productQueryRepository) FindAllProducts(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProducts(ctx, reqDb)

	if err != nil {
		return nil, nil, product_errors.ErrFindAllProducts
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordPagination(res), &totalCount, nil
}

func (r *productQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsActive(ctx, reqDb)

	if err != nil {
		return nil, nil, product_errors.ErrFindByActive
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordActivePagination(res), &totalCount, nil
}

func (r *productQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, product_errors.ErrFindByTrashed
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordTrashedPagination(res), &totalCount, nil
}

func (r *productQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    sql.NullString{String: req.Search, Valid: true},
		Column3:    req.CategoryID,
		Column4:    req.MinPrice,
		Column5:    req.MaxPrice,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetProductsByMerchant(ctx, reqDb)

	if err != nil {
		return nil, nil, product_errors.ErrFindByMerchant
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordMerchantPagination(res), &totalCount, nil
}

func (r *productQueryRepository) FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*record.ProductRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetProductsByCategoryNameParams{
		Name:    req.CategoryName,
		Column2: req.Search,
		Column3: req.MinPrice,
		Column4: req.MaxPrice,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetProductsByCategoryName(ctx, reqDb)

	if err != nil {
		return nil, nil, product_errors.ErrFindByCategory
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToProductsRecordCategoryPagination(res), &totalCount, nil
}

func (r *productQueryRepository) FindById(ctx context.Context, product_id int) (*record.ProductRecord, error) {
	res, err := r.db.GetProductByID(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrFindById
	}

	return r.mapping.ToProductRecord(res), nil
}

func (r *productQueryRepository) FindByIdTrashed(ctx context.Context, product_id int) (*record.ProductRecord, error) {
	res, err := r.db.GetProductByIdTrashed(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrFindByIdTrashed
	}

	return r.mapping.ToProductRecord(res), nil
}
