package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	refreshtoken_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/refresh_token_errors"

	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type refreshTokenRepository struct {
	db      *db.Queries
	mapping recordmapper.RefreshTokenRecordMapping
}

func NewRefreshTokenRepository(db *db.Queries, mapping recordmapper.RefreshTokenRecordMapping) *refreshTokenRepository {
	return &refreshTokenRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*record.RefreshTokenRecord, error) {
	res, err := r.db.FindRefreshTokenByToken(ctx, token)

	if err != nil {
		return nil, refreshtoken_errors.ErrTokenNotFound
	}

	return r.mapping.ToRefreshTokenRecord(res), nil
}

func (r *refreshTokenRepository) FindByUserId(ctx context.Context, user_id int) (*record.RefreshTokenRecord, error) {
	res, err := r.db.FindRefreshTokenByUserId(ctx, int32(user_id))

	if err != nil {
		return nil, refreshtoken_errors.ErrFindByUserID
	}

	return r.mapping.ToRefreshTokenRecord(res), nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*record.RefreshTokenRecord, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}

	res, err := r.db.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	})

	if err != nil {
		return nil, refreshtoken_errors.ErrCreateRefreshToken
	}

	return r.mapping.ToRefreshTokenRecord(res), nil
}

func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*record.RefreshTokenRecord, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}

	res, err := r.db.UpdateRefreshTokenByUserId(ctx, db.UpdateRefreshTokenByUserIdParams{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	})
	if err != nil {
		return nil, refreshtoken_errors.ErrUpdateRefreshToken
	}

	return r.mapping.ToRefreshTokenRecord(res), nil
}
func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	err := r.db.DeleteRefreshToken(ctx, token)

	if err != nil {
		return refreshtoken_errors.ErrDeleteRefreshToken
	}

	return nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error {
	err := r.db.DeleteRefreshTokenByUserId(ctx, int32(user_id))

	if err != nil {
		return refreshtoken_errors.ErrDeleteByUserID
	}

	return nil
}
