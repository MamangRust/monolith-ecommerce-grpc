package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantdocument_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_document_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantDocumentCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.MerchantDocumentMapping
}

func NewMerchantDocumentCommandRepository(db *db.Queries, mapping recordmapper.MerchantDocumentMapping) *merchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *merchantDocumentCommandRepository) CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error) {
	req := db.CreateMerchantDocumentParams{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         sql.NullString{String: "", Valid: true},
	}

	res, err := r.db.CreateMerchantDocument(ctx, req)
	if err != nil {
		return nil, merchantdocument_errors.ErrCreateMerchantDocumentFailed
	}

	return r.mapping.ToGetMerchantDocument(res), nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error) {
	req := db.UpdateMerchantDocumentParams{
		DocumentID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       request.Status,
		Note:         sql.NullString{String: request.Note, Valid: true},
	}

	res, err := r.db.UpdateMerchantDocument(ctx, req)
	if err != nil {
		return nil, merchantdocument_errors.ErrUpdateMerchantDocumentFailed
	}

	return r.mapping.ToGetMerchantDocument(res), nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*record.MerchantDocumentRecord, error) {
	req := db.UpdateMerchantDocumentStatusParams{
		DocumentID: int32(request.MerchantID),
		Status:     request.Status,
		Note:       sql.NullString{String: request.Note, Valid: true},
	}

	res, err := r.db.UpdateMerchantDocumentStatus(ctx, req)
	if err != nil {
		return nil, merchantdocument_errors.ErrUpdateMerchantDocumentStatusFailed
	}

	return r.mapping.ToGetMerchantDocument(res), nil
}

func (r *merchantDocumentCommandRepository) TrashedMerchantDocument(ctx context.Context, documentID int) (*record.MerchantDocumentRecord, error) {
	res, err := r.db.TrashMerchantDocument(ctx, int32(documentID))
	if err != nil {
		return nil, merchantdocument_errors.ErrTrashedMerchantDocumentFailed
	}

	return r.mapping.ToGetMerchantDocument(res), nil
}

func (r *merchantDocumentCommandRepository) RestoreMerchantDocument(ctx context.Context, documentID int) (*record.MerchantDocumentRecord, error) {
	res, err := r.db.RestoreMerchantDocument(ctx, int32(documentID))
	if err != nil {
		return nil, merchantdocument_errors.ErrRestoreMerchantDocumentFailed
	}

	return r.mapping.ToGetMerchantDocument(res), nil
}

func (r *merchantDocumentCommandRepository) DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error) {
	err := r.db.DeleteMerchantDocumentPermanently(ctx, int32(documentID))
	if err != nil {
		return false, merchantdocument_errors.ErrDeleteMerchantDocumentPermanentFailed
	}

	return true, nil
}

func (r *merchantDocumentCommandRepository) RestoreAllMerchantDocument(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantDocuments(ctx)
	if err != nil {
		return false, merchantdocument_errors.ErrRestoreAllMerchantDocumentsFailed
	}

	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDocuments(ctx)
	if err != nil {
		return false, merchantdocument_errors.ErrDeleteAllMerchantDocumentsPermanentFailed
	}

	return true, nil
}
