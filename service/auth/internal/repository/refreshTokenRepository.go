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
	ctx     context.Context
	mapping recordmapper.RefreshTokenRecordMapping
}

func NewRefreshTokenRepository(db *db.Queries, ctx context.Context, mapping recordmapper.RefreshTokenRecordMapping) *refreshTokenRepository {
	return &refreshTokenRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *refreshTokenRepository) FindByToken(token string) (*record.RefreshTokenRecord, error) {
	res, err := r.db.FindRefreshTokenByToken(r.ctx, token)

	if err != nil {
		return nil, refreshtoken_errors.ErrTokenNotFound
	}

	return r.mapping.ToRefreshTokenRecord(res), nil
}

func (r *refreshTokenRepository) FindByUserId(user_id int) (*record.RefreshTokenRecord, error) {
	res, err := r.db.FindRefreshTokenByUserId(r.ctx, int32(user_id))

	if err != nil {
		return nil, refreshtoken_errors.ErrFindByUserID
	}

	return r.mapping.ToRefreshTokenRecord(res), nil
}

func (r *refreshTokenRepository) CreateRefreshToken(req *requests.CreateRefreshToken) (*record.RefreshTokenRecord, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}

	res, err := r.db.CreateRefreshToken(r.ctx, db.CreateRefreshTokenParams{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	})

	if err != nil {
		return nil, refreshtoken_errors.ErrCreateRefreshToken
	}

	return r.mapping.ToRefreshTokenRecord(res), nil
}

func (r *refreshTokenRepository) UpdateRefreshToken(req *requests.UpdateRefreshToken) (*record.RefreshTokenRecord, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}

	res, err := r.db.UpdateRefreshTokenByUserId(r.ctx, db.UpdateRefreshTokenByUserIdParams{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	})
	if err != nil {
		return nil, refreshtoken_errors.ErrUpdateRefreshToken
	}

	return r.mapping.ToRefreshTokenRecord(res), nil
}
func (r *refreshTokenRepository) DeleteRefreshToken(token string) error {
	err := r.db.DeleteRefreshToken(r.ctx, token)

	if err != nil {
		return refreshtoken_errors.ErrDeleteRefreshToken
	}

	return nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(user_id int) error {
	err := r.db.DeleteRefreshTokenByUserId(r.ctx, int32(user_id))

	if err != nil {
		return refreshtoken_errors.ErrDeleteByUserID
	}

	return nil
}
