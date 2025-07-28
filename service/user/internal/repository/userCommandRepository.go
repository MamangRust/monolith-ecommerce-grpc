package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type userCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.UserRecordMapping
}

func NewUserCommandRepository(db *db.Queries, mapping recordmapper.UserRecordMapping) *userCommandRepository {
	return &userCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *userCommandRepository) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*record.UserRecord, error) {
	req := db.CreateUserParams{
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	user, err := r.db.CreateUser(ctx, req)

	if err != nil {
		return nil, user_errors.ErrCreateUser
	}

	return r.mapping.ToUserRecord(user), nil
}

func (r *userCommandRepository) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*record.UserRecord, error) {
	req := db.UpdateUserParams{
		UserID:    int32(*request.UserID),
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	res, err := r.db.UpdateUser(ctx, req)

	if err != nil {
		return nil, user_errors.ErrUpdateUser
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userCommandRepository) TrashedUser(ctx context.Context, user_id int) (*record.UserRecord, error) {
	res, err := r.db.TrashUser(ctx, int32(user_id))

	if err != nil {
		return nil, user_errors.ErrTrashedUser
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userCommandRepository) RestoreUser(ctx context.Context, user_id int) (*record.UserRecord, error) {
	res, err := r.db.RestoreUser(ctx, int32(user_id))

	if err != nil {
		return nil, user_errors.ErrRestoreUser
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userCommandRepository) DeleteUserPermanent(ctx context.Context, user_id int) (bool, error) {
	err := r.db.DeleteUserPermanently(ctx, int32(user_id))

	if err != nil {
		return false, user_errors.ErrDeleteUserPermanent
	}

	return true, nil
}

func (r *userCommandRepository) RestoreAllUser(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllUsers(ctx)

	if err != nil {
		return false, user_errors.ErrRestoreAllUsers
	}

	return true, nil
}

func (r *userCommandRepository) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentUsers(ctx)

	if err != nil {
		return false, user_errors.ErrDeleteAllUsers
	}
	return true, nil
}
