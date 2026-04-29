package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
)


type roleQueryRepository struct {
	db *db.Queries
}

func NewRoleQueryRepository(db *db.Queries) *roleQueryRepository {
	return &roleQueryRepository{
		db: db,
	}
}

func (r *roleQueryRepository) FindAll(ctx context.Context, req *requests.FindAllRole) ([]*db.GetRolesRow, error) {
	fmt.Printf("DEBUG: FindAllRoles search='%s'\n", req.Search)
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetRoles(ctx, reqDb)

	if err != nil {
		fmt.Printf("DEBUG: FindAllRoles db error: %v\n", err)
		return nil, role_errors.ErrFindAllRoles.WithInternal(err)
	}

	fmt.Printf("DEBUG: FindAllRoles found %d roles\n", len(res))
	if len(res) > 0 {
		fmt.Printf("DEBUG: Role[0] Name: %s\n", res[0].RoleName)
	}

	return res, nil
}

func (r *roleQueryRepository) FindByID(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.GetRole(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}


func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*db.Role, error) {
	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}


func (r *roleQueryRepository) FindByUserId(ctx context.Context, user_id int) ([]*db.Role, error) {
	res, err := r.db.GetUserRoles(ctx, int32(user_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}


func (r *roleQueryRepository) FindActive(ctx context.Context, req *requests.FindAllRole) ([]*db.GetActiveRolesRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveRoles(ctx, reqDb)

	if err != nil {
		return nil, role_errors.ErrFindActiveRoles.WithInternal(err)
	}


	return res, nil
}

func (r *roleQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllRole) ([]*db.GetTrashedRolesRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedRoles(ctx, reqDb)

	if err != nil {
		return nil, role_errors.ErrFindTrashedRoles.WithInternal(err)
	}


	return res, nil
}
