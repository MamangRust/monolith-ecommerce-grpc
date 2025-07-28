package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type resetTokenRepository struct {
	db      *db.Queries
	mapping recordmapper.ResetTokenRecordMapping
}

func NewResetTokenRepository(db *db.Queries, mapping recordmapper.ResetTokenRecordMapping) *resetTokenRepository {
	return &resetTokenRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *resetTokenRepository) FindByToken(ctx context.Context, code string) (*record.ResetTokenRecord, error) {
	res, err := r.db.GetResetToken(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.mapping.ToResetTokenRecord(res), nil
}

func (r *resetTokenRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*record.ResetTokenRecord, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, err
	}
	res, err := r.db.CreateResetToken(ctx, db.CreateResetTokenParams{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	})
	if err != nil {
		return nil, err
	}
	return r.mapping.ToResetTokenRecord(res), nil
}
func (r *resetTokenRepository) DeleteResetToken(ctx context.Context, user_id int) error {
	err := r.db.DeleteResetToken(ctx, int64(user_id))
	if err != nil {
		return err
	}
	return nil
}
