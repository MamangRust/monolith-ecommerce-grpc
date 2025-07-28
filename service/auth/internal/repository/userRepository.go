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
	mapping recordmapper.UserRecordMapping
}

func NewUserRepository(db *db.Queries, mapping recordmapper.UserRecordMapping) *userRepository {
	return &userRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *userRepository) FindById(ctx context.Context, user_id int) (*record.UserRecord, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, user_errors.ErrUserNotFound
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) FindByEmailAndVerify(ctx context.Context, email string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByEmailAndVerify(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) FindByVerificationCode(ctx context.Context, verification_code string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByVerificationCode(ctx, verification_code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound

	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*record.UserRecord, error) {
	req := db.CreateUserParams{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       sql.NullBool{Bool: request.IsVerified, Valid: true},
	}

	user, err := r.db.CreateUser(ctx, req)

	if err != nil {
		return nil, user_errors.ErrCreateUser
	}

	return r.mapping.ToUserRecord(user), nil
}

func (r *userRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*record.UserRecord, error) {
	res, err := r.db.UpdateUserIsVerified(ctx, db.UpdateUserIsVerifiedParams{
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

func (r *userRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*record.UserRecord, error) {
	res, err := r.db.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		UserID:   int32(user_id),
		Password: password,
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword
	}

	return r.mapping.ToUserRecord(res), nil
}
