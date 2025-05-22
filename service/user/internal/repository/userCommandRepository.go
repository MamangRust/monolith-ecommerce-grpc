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
	ctx     context.Context
	mapping recordmapper.UserRecordMapping
}

func NewUserCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.UserRecordMapping) *userCommandRepository {
	return &userCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *userCommandRepository) CreateUser(request *requests.CreateUserRequest) (*record.UserRecord, error) {
	req := db.CreateUserParams{
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	user, err := r.db.CreateUser(r.ctx, req)

	if err != nil {
		return nil, user_errors.ErrCreateUser
	}

	return r.mapping.ToUserRecord(user), nil
}

func (r *userCommandRepository) UpdateUser(request *requests.UpdateUserRequest) (*record.UserRecord, error) {
	req := db.UpdateUserParams{
		UserID:    int32(*request.UserID),
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	res, err := r.db.UpdateUser(r.ctx, req)

	if err != nil {
		return nil, user_errors.ErrUpdateUser
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userCommandRepository) TrashedUser(user_id int) (*record.UserRecord, error) {
	res, err := r.db.TrashUser(r.ctx, int32(user_id))

	if err != nil {
		return nil, user_errors.ErrTrashedUser
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userCommandRepository) RestoreUser(user_id int) (*record.UserRecord, error) {
	res, err := r.db.RestoreUser(r.ctx, int32(user_id))

	if err != nil {
		return nil, user_errors.ErrRestoreUser
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userCommandRepository) DeleteUserPermanent(user_id int) (bool, error) {
	err := r.db.DeleteUserPermanently(r.ctx, int32(user_id))

	if err != nil {
		return false, user_errors.ErrDeleteUserPermanent
	}

	return true, nil
}

func (r *userCommandRepository) RestoreAllUser() (bool, error) {
	err := r.db.RestoreAllUsers(r.ctx)

	if err != nil {
		return false, user_errors.ErrRestoreAllUsers
	}

	return true, nil
}

func (r *userCommandRepository) DeleteAllUserPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentUsers(r.ctx)

	if err != nil {
		return false, user_errors.ErrDeleteAllUsers
	}
	return true, nil
}
