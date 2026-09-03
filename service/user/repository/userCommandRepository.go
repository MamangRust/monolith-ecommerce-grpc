package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
)

type userCommandRepository struct {
	db *db.Queries
}

func NewUserCommandRepository(db *db.Queries) *userCommandRepository {
	return &userCommandRepository{
		db: db,
	}
}

func (r *userCommandRepository) Create(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error) {
	req := db.CreateUserParams{
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	user, err := r.db.CreateUser(ctx, req)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}

	return user, nil
}

func (r *userCommandRepository) Update(ctx context.Context, request *requests.UpdateUserRequest) (*db.User, error) {
	req := db.UpdateUserParams{
		UserID:    int32(*request.UserID),
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	res, err := r.db.UpdateUser(ctx, req)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}

	return res, nil
}

func (r *userCommandRepository) Trash(ctx context.Context, user_id int) (*db.TrashUserRow, error) {
	res, err := r.db.TrashUser(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrTrashedUser.WithInternal(err)
	}

	return res, nil
}

func (r *userCommandRepository) Restore(ctx context.Context, user_id int) (*db.RestoreUserRow, error) {
	res, err := r.db.RestoreUser(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrRestoreUser.WithInternal(err)
	}

	return res, nil
}

func (r *userCommandRepository) DeletePermanent(ctx context.Context, user_id int) (bool, error) {
	err := r.db.DeleteUserPermanently(ctx, int32(user_id))

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return false, shared_errors.NewConflictError("cannot permanently delete user while related records exist").WithInternal(err)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return false, user_errors.ErrUserNotFound
		}
		return false, user_errors.ErrDeleteUserPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *userCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllUsers(ctx)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, user_errors.ErrUserNotFound
		}
		return false, user_errors.ErrRestoreAllUsers.WithInternal(err)
	}

	return true, nil
}

func (r *userCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentUsers(ctx)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return false, shared_errors.NewConflictError("cannot permanently delete users while related records exist").WithInternal(err)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return false, user_errors.ErrUserNotFound
		}
		return false, user_errors.ErrDeleteAllUsers.WithInternal(err)
	}
	return true, nil
}

func (r *userCommandRepository) UpdateIsVerified(ctx context.Context, user_id int, is_verified bool) (*db.UpdateUserIsVerifiedRow, error) {
	arg := db.UpdateUserIsVerifiedParams{
		UserID:     int32(user_id),
		IsVerified: &is_verified,
	}

	res, err := r.db.UpdateUserIsVerified(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}

	return res, nil
}

func (r *userCommandRepository) UpdatePassword(ctx context.Context, user_id int, password string) (*db.UpdateUserPasswordRow, error) {
	arg := db.UpdateUserPasswordParams{
		UserID:   int32(user_id),
		Password: password,
	}

	res, err := r.db.UpdateUserPassword(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}

	return res, nil
}
