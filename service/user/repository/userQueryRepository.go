package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
)



type userQueryRepository struct {
	db *db.Queries
}

func NewUserQueryRepository(db *db.Queries) *userQueryRepository {
	return &userQueryRepository{
		db: db,
	}
}

func (r *userQueryRepository) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsers(ctx, reqDb)

	if err != nil {
		return nil, user_errors.ErrFindAllUsers.WithInternal(err)
	}


	return res, nil
}

func (r *userQueryRepository) FindByID(ctx context.Context, user_id int) (*db.GetUserByIDRow, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}


	return res, nil
}


func (r *userQueryRepository) FindByIDWithPassword(ctx context.Context, user_id int) (*db.GetUserByIDRow, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}


	return res, nil
}


func (r *userQueryRepository) FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsersActive(ctx, reqDb)

	if err != nil {
		return nil, user_errors.ErrFindActiveUsers.WithInternal(err)
	}


	return res, nil
}

func (r *userQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUserTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUserTrashed(ctx, reqDb)

	if err != nil {
		return nil, user_errors.ErrFindTrashedUsers.WithInternal(err)
	}


	return res, nil
}

func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*db.User, error) {
	res, err := r.db.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}


	return res, nil
}


func (r *userQueryRepository) FindByEmailWithPassword(ctx context.Context, email string) (*db.GetUserByEmailWithPasswordRow, error) {
	res, err := r.db.GetUserByEmailWithPassword(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}


	return res, nil
}

func (r *userQueryRepository) FindByVerificationCode(ctx context.Context, code string) (*db.GetUserByVerificationCodeRow, error) {
	res, err := r.db.GetUserByVerificationCode(ctx, code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}
