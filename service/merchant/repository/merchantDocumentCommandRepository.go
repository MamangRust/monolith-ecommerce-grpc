package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
)


type merchantDocumentCommandRepository struct {
	db *db.Queries
}

func NewMerchantDocumentCommandRepository(db *db.Queries) *merchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{
		db: db,
	}
}

func (r *merchantDocumentCommandRepository) Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.CreateMerchantDocumentRow, error) {
	req := db.CreateMerchantDocumentParams{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         stringPtr(""),
	}

	res, err := r.db.CreateMerchantDocument(ctx, req)
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	return res, nil
}

func (r *merchantDocumentCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.UpdateMerchantDocumentRow, error) {
	req := db.UpdateMerchantDocumentParams{
		DocumentID:   int32(*request.DocumentID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       request.Status,
		Note:         stringPtr(request.Note),
	}

	res, err := r.db.UpdateMerchantDocument(ctx, req)
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	return res, nil
}

func (r *merchantDocumentCommandRepository) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error) {
	req := db.UpdateMerchantDocumentStatusParams{
		DocumentID: int32(*request.DocumentID),
		Status:     request.Status,
		Note:       stringPtr(request.Note),
	}

	res, err := r.db.UpdateMerchantDocumentStatus(ctx, req)
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	return res, nil
}

func (r *merchantDocumentCommandRepository) Trash(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	res, err := r.db.TrashMerchantDocument(ctx, int32(documentID))
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	return res, nil
}

func (r *merchantDocumentCommandRepository) Restore(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	res, err := r.db.RestoreMerchantDocument(ctx, int32(documentID))
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	return res, nil
}

func (r *merchantDocumentCommandRepository) DeletePermanent(ctx context.Context, documentID int) (bool, error) {
	err := r.db.DeleteMerchantDocumentPermanently(ctx, int32(documentID))
	if err != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	return true, nil
}

func (r *merchantDocumentCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantDocuments(ctx)
	if err != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDocuments(ctx)
	if err != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	return true, nil
}
