package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type roleCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.RoleRecordMapping
}

func NewRoleCommandRepository(db *db.Queries, mapping recordmapper.RoleRecordMapping) *roleCommandRepository {
	return &roleCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *roleCommandRepository) CreateRole(ctx context.Context, req *requests.CreateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.CreateRole(ctx, req.Name)

	if err != nil {
		return nil, role_errors.ErrCreateRole
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) UpdateRole(ctx context.Context, req *requests.UpdateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.UpdateRole(ctx, db.UpdateRoleParams{
		RoleID:   int32(*req.ID),
		RoleName: req.Name,
	})

	if err != nil {
		return nil, role_errors.ErrUpdateRole
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) TrashedRole(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.TrashRole(ctx, int32(id))
	if err != nil {
		return nil, role_errors.ErrTrashedRole
	}
	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) RestoreRole(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.RestoreRole(ctx, int32(id))
	if err != nil {
		return nil, role_errors.ErrRestoreRole
	}
	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) DeleteRolePermanent(ctx context.Context, role_id int) (bool, error) {
	err := r.db.DeletePermanentRole(ctx, int32(role_id))
	if err != nil {
		return false, role_errors.ErrDeleteRolePermanent
	}
	return true, nil
}

func (r *roleCommandRepository) RestoreAllRole(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllRoles(ctx)

	if err != nil {
		return false, role_errors.ErrRestoreAllRoles
	}

	return true, nil
}

func (r *roleCommandRepository) DeleteAllRolePermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentRoles(ctx)

	if err != nil {
		return false, role_errors.ErrDeleteAllRoles
	}

	return true, nil
}
