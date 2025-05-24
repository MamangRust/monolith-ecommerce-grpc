package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	merchantdocument_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_document_errors"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type merchantDocumentCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantDocumentMapping
}

func NewMerchantDocumentCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantDocumentMapping) *merchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantDocumentCommandRepository) CreateMerchantDocument(request *requests.CreateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error) {
	req := db.CreateMerchantDocumentParams{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         sql.NullString{String: "", Valid: true},
	}

	res, err := r.db.CreateMerchantDocument(r.ctx, req)
	if err != nil {
		return nil, merchantdocument_errors.ErrCreateMerchantDocumentFailed
	}

	return r.mapping.ToGetMerchantDocument(res), nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocument(request *requests.UpdateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error) {
	req := db.UpdateMerchantDocumentParams{
		DocumentID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       request.Status,
		Note:         sql.NullString{String: request.Note, Valid: true},
	}

	res, err := r.db.UpdateMerchantDocument(r.ctx, req)
	if err != nil {
		return nil, merchantdocument_errors.ErrUpdateMerchantDocumentFailed
	}

	return r.mapping.ToGetMerchantDocument(res), nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatus(request *requests.UpdateMerchantDocumentStatusRequest) (*record.MerchantDocumentRecord, error) {
	req := db.UpdateMerchantDocumentStatusParams{
		DocumentID: int32(request.MerchantID),
		Status:     request.Status,
		Note:       sql.NullString{String: request.Note, Valid: true},
	}

	res, err := r.db.UpdateMerchantDocumentStatus(r.ctx, req)
	if err != nil {
		return nil, merchantdocument_errors.ErrUpdateMerchantDocumentStatusFailed
	}

	return r.mapping.ToGetMerchantDocument(res), nil
}

func (r *merchantDocumentCommandRepository) TrashedMerchantDocument(documentID int) (*record.MerchantDocumentRecord, error) {
	res, err := r.db.TrashMerchantDocument(r.ctx, int32(documentID))
	if err != nil {
		return nil, merchantdocument_errors.ErrTrashedMerchantDocumentFailed
	}

	return r.mapping.ToGetMerchantDocument(res), nil
}

func (r *merchantDocumentCommandRepository) RestoreMerchantDocument(documentID int) (*record.MerchantDocumentRecord, error) {
	res, err := r.db.RestoreMerchantDocument(r.ctx, int32(documentID))
	if err != nil {
		return nil, merchantdocument_errors.ErrRestoreMerchantDocumentFailed
	}

	return r.mapping.ToGetMerchantDocument(res), nil
}

func (r *merchantDocumentCommandRepository) DeleteMerchantDocumentPermanent(documentID int) (bool, error) {
	err := r.db.DeleteMerchantDocumentPermanently(r.ctx, int32(documentID))
	if err != nil {
		return false, merchantdocument_errors.ErrDeleteMerchantDocumentPermanentFailed
	}

	return true, nil
}

func (r *merchantDocumentCommandRepository) RestoreAllMerchantDocument() (bool, error) {
	err := r.db.RestoreAllMerchantDocuments(r.ctx)
	if err != nil {
		return false, merchantdocument_errors.ErrRestoreAllMerchantDocumentsFailed
	}

	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAllMerchantDocumentPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDocuments(r.ctx)
	if err != nil {
		return false, merchantdocument_errors.ErrDeleteAllMerchantDocumentsPermanentFailed
	}

	return true, nil
}
