package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/convert"
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
		Note:         convert.NullableString(""),
	}

	res, err := r.db.CreateMerchantDocument(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return res, nil
}

// CreateInTx persists the document inside the given database transaction so the
// caller can commit the business write and its outbox event atomically (Phase 6
// — transactional outbox).
func (r *merchantDocumentCommandRepository) CreateInTx(ctx context.Context, tx pgx.Tx, request *requests.CreateMerchantDocumentRequest) (*db.CreateMerchantDocumentRow, error) {
	req := db.CreateMerchantDocumentParams{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         convert.NullableString(""),
	}

	res, err := r.db.WithTx(tx).CreateMerchantDocument(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
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
		Note:         convert.NullableString(request.Note),
	}

	res, err := r.db.UpdateMerchantDocument(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDocumentCommandRepository) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error) {
	req := db.UpdateMerchantDocumentStatusParams{
		DocumentID: int32(*request.DocumentID),
		Status:     request.Status,
		Note:       convert.NullableString(request.Note),
	}

	res, err := r.db.UpdateMerchantDocumentStatus(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return res, nil
}

// UpdateStatusInTx updates the document status inside the given database
// transaction so the caller can commit the business write and its outbox event
// atomically (Phase 6 — transactional outbox).
func (r *merchantDocumentCommandRepository) UpdateStatusInTx(ctx context.Context, tx pgx.Tx, request *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error) {
	req := db.UpdateMerchantDocumentStatusParams{
		DocumentID: int32(*request.DocumentID),
		Status:     request.Status,
		Note:       convert.NullableString(request.Note),
	}

	res, err := r.db.WithTx(tx).UpdateMerchantDocumentStatus(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDocumentCommandRepository) Trash(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	res, err := r.db.TrashMerchantDocument(ctx, int32(documentID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDocumentCommandRepository) Restore(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	res, err := r.db.RestoreMerchantDocument(ctx, int32(documentID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDocumentCommandRepository) DeletePermanent(ctx context.Context, documentID int) (bool, error) {
	err := r.db.DeleteMerchantDocumentPermanently(ctx, int32(documentID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, merchant_errors.ErrMerchantNotFound
		}
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return true, nil
}

func (r *merchantDocumentCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantDocuments(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, merchant_errors.ErrMerchantNotFound
		}
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDocuments(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, merchant_errors.ErrMerchantNotFound
		}
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return true, nil
}
