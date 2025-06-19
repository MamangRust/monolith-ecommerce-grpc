package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type userRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.UserRecordMapping
}

func NewUserRepository(db *db.Queries, ctx context.Context, mapping recordmapper.UserRecordMapping) *userRepository {
	return &userRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *userRepository) FindById(user_id int) (*record.UserRecord, error) {
	res, err := r.db.GetUserByID(r.ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) FindByEmail(email string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByEmail(r.ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) FindByEmailAndVerify(email string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByEmailAndVerify(r.ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) FindByVerificationCode(verification_code string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByVerificationCode(r.ctx, verification_code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound

	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) CreateUser(request *requests.RegisterRequest) (*record.UserRecord, error) {
	req := db.CreateUserParams{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       sql.NullBool{Bool: request.IsVerified, Valid: true},
	}

	user, err := r.db.CreateUser(r.ctx, req)

	if err != nil {
		return nil, user_errors.ErrCreateUser
	}

	return r.mapping.ToUserRecord(user), nil
}

func (r *userRepository) UpdateUserIsVerified(user_id int, is_verified bool) (*record.UserRecord, error) {
	res, err := r.db.UpdateUserIsVerified(r.ctx, db.UpdateUserIsVerifiedParams{
		UserID: int32(user_id),
		IsVerified: sql.NullBool{
			Bool:  is_verified,
			Valid: true,
		},
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) UpdateUserPassword(user_id int, password string) (*record.UserRecord, error) {
	res, err := r.db.UpdateUserPassword(r.ctx, db.UpdateUserPasswordParams{
		UserID:   int32(user_id),
		Password: password,
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword
	}

	return r.mapping.ToUserRecord(res), nil
}
