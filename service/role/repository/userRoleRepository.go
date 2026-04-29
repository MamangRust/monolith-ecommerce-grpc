package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
)

type userRoleRepository struct {
	db *db.Queries
}

func NewUserRoleRepository(db *db.Queries) UserRoleRepository {
	return &userRoleRepository{
		db: db,
	}
}

func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error) {
	arg := db.AssignRoleToUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	}

	res, err := r.db.AssignRoleToUser(ctx, arg)
	if err != nil {
		return nil, role_errors.ErrAssignRole.WithInternal(err)
	}

	return res, nil
}

func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	arg := db.RemoveRoleFromUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	}

	err := r.db.RemoveRoleFromUser(ctx, arg)
	if err != nil {
		return role_errors.ErrRemoveRole.WithInternal(err)
	}

	return nil
}
