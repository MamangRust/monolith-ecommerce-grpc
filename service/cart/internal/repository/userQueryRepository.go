package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
)

type userQueryRepository struct {
	db *db.Queries
}

func NewUserQueryRepository(db *db.Queries) UserQueryRepository {
	return &userQueryRepository{
		db: db,
	}
}


func (r *userQueryRepository) FindById(ctx context.Context, user_id int) (*db.GetUserByIDRow, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, user_errors.ErrUserInternal.WithInternal(err)
	}

	return res, nil
}

