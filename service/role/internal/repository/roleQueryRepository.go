package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type roleQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.RoleRecordMapping
}

func NewRoleQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.RoleRecordMapping) *roleQueryRepository {
	return &roleQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *roleQueryRepository) FindAllRoles(req *requests.FindAllRole) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetRoles(r.ctx, reqDb)

	if err != nil {
		return nil, nil, role_errors.ErrFindAllRoles
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToRolesRecordAll(res), &totalCount, nil
}

func (r *roleQueryRepository) FindById(id int) (*record.RoleRecord, error) {
	res, err := r.db.GetRole(r.ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found with ID: %d", id)
		}
		return nil, fmt.Errorf("failed to find role by ID %d: %w", id, err)
	}
	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleQueryRepository) FindByName(name string) (*record.RoleRecord, error) {
	res, err := r.db.GetRoleByName(r.ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound
		}

		return nil, role_errors.ErrRoleNotFound
	}
	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleQueryRepository) FindByUserId(user_id int) ([]*record.RoleRecord, error) {
	res, err := r.db.GetUserRoles(r.ctx, int32(user_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound
		}

		return nil, role_errors.ErrRoleNotFound
	}
	return r.mapping.ToRolesRecord(res), nil
}

func (r *roleQueryRepository) FindByActiveRole(req *requests.FindAllRole) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveRoles(r.ctx, reqDb)

	if err != nil {
		return nil, nil, role_errors.ErrFindActiveRoles
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToRolesRecordActive(res), &totalCount, nil
}

func (r *roleQueryRepository) FindByTrashedRole(req *requests.FindAllRole) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedRoles(r.ctx, reqDb)

	if err != nil {
		return nil, nil, role_errors.ErrFindTrashedRoles
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToRolesRecordTrashed(res), &totalCount, nil
}
