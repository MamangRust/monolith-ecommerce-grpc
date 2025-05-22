package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type userQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.UserRecordMapping
}

func NewUserQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.UserRecordMapping) *userQueryRepository {
	return &userQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *userQueryRepository) FindById(user_id int) (*record.UserRecord, error) {
	res, err := r.db.GetUserByID(r.ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	return r.mapping.ToUserRecord(res), nil
}
