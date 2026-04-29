package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
)


type roleCommandRepository struct {
	db *db.Queries
}

func NewRoleCommandRepository(db *db.Queries) *roleCommandRepository {
	return &roleCommandRepository{
		db: db,
	}
}

func (r *roleCommandRepository) Create(ctx context.Context, req *requests.CreateRoleRequest) (*db.Role, error) {
	res, err := r.db.CreateRole(ctx, req.Name)

	if err != nil {
		return nil, role_errors.ErrCreateRole.WithInternal(err)
	}


	return res, nil
}

func (r *roleCommandRepository) Update(ctx context.Context, req *requests.UpdateRoleRequest) (*db.Role, error) {
	res, err := r.db.UpdateRole(ctx, db.UpdateRoleParams{
		RoleID:   int32(*req.ID),
		RoleName: req.Name,
	})

	if err != nil {
		return nil, role_errors.ErrUpdateRole.WithInternal(err)
	}


	return res, nil
}

func (r *roleCommandRepository) Trash(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.TrashRole(ctx, int32(id))
	if err != nil {
		return nil, role_errors.ErrTrashedRole.WithInternal(err)
	}
	return res, nil
}


func (r *roleCommandRepository) Restore(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.RestoreRole(ctx, int32(id))
	if err != nil {
		return nil, role_errors.ErrRestoreRole.WithInternal(err)
	}
	return res, nil
}


func (r *roleCommandRepository) DeletePermanent(ctx context.Context, role_id int) (bool, error) {
	err := r.db.DeletePermanentRole(ctx, int32(role_id))
	if err != nil {
		return false, role_errors.ErrDeleteRolePermanent.WithInternal(err)
	}
	return true, nil
}


func (r *roleCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllRoles(ctx)

	if err != nil {
		return false, role_errors.ErrRestoreAllRoles.WithInternal(err)
	}


	return true, nil
}

func (r *roleCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentRoles(ctx)

	if err != nil {
		return false, role_errors.ErrDeleteAllRoles.WithInternal(err)
	}


	return true, nil
}
