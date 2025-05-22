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
	ctx     context.Context
	mapping recordmapper.RoleRecordMapping
}

func NewRoleCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.RoleRecordMapping) *roleCommandRepository {
	return &roleCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *roleCommandRepository) CreateRole(req *requests.CreateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.CreateRole(r.ctx, req.Name)

	if err != nil {
		return nil, role_errors.ErrCreateRole
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) UpdateRole(req *requests.UpdateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.UpdateRole(r.ctx, db.UpdateRoleParams{
		RoleID:   int32(*req.ID),
		RoleName: req.Name,
	})

	if err != nil {
		return nil, role_errors.ErrUpdateRole
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) TrashedRole(id int) (*record.RoleRecord, error) {
	res, err := r.db.TrashRole(r.ctx, int32(id))
	if err != nil {
		return nil, role_errors.ErrTrashedRole
	}
	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) RestoreRole(id int) (*record.RoleRecord, error) {
	res, err := r.db.RestoreRole(r.ctx, int32(id))
	if err != nil {
		return nil, role_errors.ErrRestoreRole
	}
	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) DeleteRolePermanent(role_id int) (bool, error) {
	err := r.db.DeletePermanentRole(r.ctx, int32(role_id))
	if err != nil {
		return false, role_errors.ErrDeleteRolePermanent
	}
	return true, nil
}

func (r *roleCommandRepository) RestoreAllRole() (bool, error) {
	err := r.db.RestoreAllRoles(r.ctx)

	if err != nil {
		return false, role_errors.ErrRestoreAllRoles
	}

	return true, nil
}

func (r *roleCommandRepository) DeleteAllRolePermanent() (bool, error) {
	err := r.db.DeleteAllPermanentRoles(r.ctx)

	if err != nil {
		return false, role_errors.ErrDeleteAllRoles
	}

	return true, nil
}
