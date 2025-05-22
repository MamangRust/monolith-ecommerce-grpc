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
	ctx     context.Context
	mapping recordmapper.ResetTokenRecordMapping
}

func NewResetTokenRepository(db *db.Queries, ctx context.Context, mapping recordmapper.ResetTokenRecordMapping) *resetTokenRepository {
	return &resetTokenRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *resetTokenRepository) FindByToken(code string) (*record.ResetTokenRecord, error) {
	res, err := r.db.GetResetToken(r.ctx, code)
	if err != nil {
		return nil, err
	}
	return r.mapping.ToResetTokenRecord(res), nil
}

func (r *resetTokenRepository) CreateResetToken(req *requests.CreateResetTokenRequest) (*record.ResetTokenRecord, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, err
	}
	res, err := r.db.CreateResetToken(r.ctx, db.CreateResetTokenParams{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	})
	if err != nil {
		return nil, err
	}
	return r.mapping.ToResetTokenRecord(res), nil
}
func (r *resetTokenRepository) DeleteResetToken(user_id int) error {
	err := r.db.DeleteResetToken(r.ctx, int64(user_id))
	if err != nil {
		return err
	}
	return nil
}
