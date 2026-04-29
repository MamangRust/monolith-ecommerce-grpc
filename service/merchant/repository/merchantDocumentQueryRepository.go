package repository

import (
	"context"

	"database/sql"
 
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
)


type merchantDocumentQueryRepository struct {
	db *db.Queries
}

func NewMerchantDocumentQueryRepository(db *db.Queries) *merchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{
		db: db,
	}
}

func (r *merchantDocumentQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetMerchantDocuments(ctx, params)
	if err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return docs, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetActiveMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetActiveMerchantDocuments(ctx, params)
	if err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return docs, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetTrashedMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetTrashedMerchantDocuments(ctx, params)
	if err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return docs, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindByID(ctx context.Context, id int) (*db.GetMerchantDocumentRow, error) {
	doc, err := r.db.GetMerchantDocument(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return doc, nil
}
