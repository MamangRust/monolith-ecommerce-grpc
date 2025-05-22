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

func (r *userQueryRepository) FindAllUsers(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsers(r.ctx, reqDb)

	if err != nil {
		return nil, nil, user_errors.ErrFindAllUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToUsersRecordPagination(res), &totalCount, nil
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

func (r *userQueryRepository) FindByActive(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsersActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, user_errors.ErrFindActiveUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToUsersRecordActivePagination(res), &totalCount, nil
}

func (r *userQueryRepository) FindByTrashed(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUserTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUserTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, user_errors.ErrFindTrashedUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToUsersRecordTrashedPagination(res), &totalCount, nil
}

func (r *userQueryRepository) FindByEmail(email string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByEmail(r.ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	return r.mapping.ToUserRecord(res), nil
}
