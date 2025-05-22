package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantdocument_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_document_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantDocumentQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantDocumentMapping
}

func NewMerchantDocumentQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantDocumentMapping) *merchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantDocumentQueryRepository) FindAllDocuments(req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetMerchantDocuments(r.ctx, params)
	if err != nil {
		return nil, nil, merchantdocument_errors.ErrFindAllMerchantDocumentsFailed
	}

	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return r.mapping.ToMerchantDocumentsRecord(docs), &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindByActive(req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetActiveMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetActiveMerchantDocuments(r.ctx, params)
	if err != nil {
		return nil, nil, merchantdocument_errors.ErrFindActiveMerchantDocumentsFailed
	}

	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return r.mapping.ToMerchantDocumentsActiveRecord(docs), &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindByTrashed(req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetTrashedMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetTrashedMerchantDocuments(r.ctx, params)
	if err != nil {
		return nil, nil, merchantdocument_errors.ErrFindTrashedMerchantDocumentsFailed
	}

	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return r.mapping.ToMerchantDocumentsTrashedRecord(docs), &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindById(id int) (*record.MerchantDocumentRecord, error) {
	doc, err := r.db.GetMerchantDocument(r.ctx, int32(id))
	if err != nil {
		return nil, merchantdocument_errors.ErrFindMerchantDocumentByIdFailed
	}
	return r.mapping.ToGetMerchantDocument(doc), nil
}
