package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/convert"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
)

type merchantCommandRepository struct {
	db *db.Queries
}

func NewMerchantCommandRepository(db *db.Queries) *merchantCommandRepository {
	return &merchantCommandRepository{
		db: db,
	}
}

func (r *merchantCommandRepository) Create(
	ctx context.Context,
	request *requests.CreateMerchantRequest,
) (*db.CreateMerchantRow, error) {
	req := db.CreateMerchantParams{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Status:       "active",
		Description:  convert.NullableString(request.Description),
		Address:      convert.NullableString(request.Address),
		ContactEmail: convert.NullableString(request.ContactEmail),
		ContactPhone: convert.NullableString(request.ContactPhone),
	}

	merchant, err := r.db.CreateMerchant(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrCreateMerchant.WithInternal(err)
	}

	return merchant, nil
}

// CreateInTx persists the merchant inside the given database transaction so the
// caller can commit the business write and its outbox event atomically (Phase 6
// — transactional outbox).
func (r *merchantCommandRepository) CreateInTx(
	ctx context.Context,
	tx pgx.Tx,
	request *requests.CreateMerchantRequest,
) (*db.CreateMerchantRow, error) {
	req := db.CreateMerchantParams{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Status:       "active",
		Description:  convert.NullableString(request.Description),
		Address:      convert.NullableString(request.Address),
		ContactEmail: convert.NullableString(request.ContactEmail),
		ContactPhone: convert.NullableString(request.ContactPhone),
	}

	merchant, err := r.db.WithTx(tx).CreateMerchant(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrCreateMerchant.WithInternal(err)
	}

	return merchant, nil
}

func (r *merchantCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error) {
	req := db.UpdateMerchantParams{
		MerchantID:   int32(*request.MerchantID),
		Name:         request.Name,
		Description:  convert.NullableString(request.Description),
		Address:      convert.NullableString(request.Address),
		ContactEmail: convert.NullableString(request.ContactEmail),
		ContactPhone: convert.NullableString(request.ContactPhone),
		Status:       request.Status,
	}

	res, err := r.db.UpdateMerchant(ctx, req)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrUpdateMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) Trash(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	res, err := r.db.TrashMerchant(ctx, int32(merchant_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrTrashedMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) Restore(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	res, err := r.db.RestoreMerchant(ctx, int32(merchant_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrRestoreMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) DeletePermanent(ctx context.Context, Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(ctx, int32(Merchant_id))

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return false, shared_errors.NewConflictError("cannot permanently delete merchant while related records exist").WithInternal(err)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return false, merchant_errors.ErrMerchantNotFound
		}
		return false, merchant_errors.ErrDeleteMerchantPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *merchantCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchants(ctx)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, merchant_errors.ErrMerchantNotFound
		}
		return false, merchant_errors.ErrRestoreAllMerchants.WithInternal(err)
	}

	return true, nil
}

func (r *merchantCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(ctx)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return false, shared_errors.NewConflictError("cannot permanently delete merchants while related records exist").WithInternal(err)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return false, merchant_errors.ErrMerchantNotFound
		}
		return false, merchant_errors.ErrDeleteAllMerchants.WithInternal(err)
	}

	return true, nil
}

func (r *merchantCommandRepository) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.UpdateMerchantStatusRow, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(*request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchantStatus(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return res, nil
}

// UpdateStatusInTx updates the merchant status inside the given database
// transaction so the caller can commit the business write and its outbox event
// atomically (Phase 6 — transactional outbox).
func (r *merchantCommandRepository) UpdateStatusInTx(ctx context.Context, tx pgx.Tx, request *requests.UpdateMerchantStatusRequest) (*db.UpdateMerchantStatusRow, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(*request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.WithTx(tx).UpdateMerchantStatus(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return res, nil
}
