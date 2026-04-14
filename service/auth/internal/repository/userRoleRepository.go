package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	userrole_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/user_role_errors"
)


// userRoleRepository is a struct that implements the UserRoleRepository interface
type userRoleRepository struct {
	db *db.Queries
}

// NewUserRoleRepository creates a new UserRoleRepository instance
//
// Args:
// db: a pointer to the database queries
// mapper: a UserRoleRecordMapping object
//
// Returns:
// a pointer to the userRoleRepository struct
func NewUserRoleRepository(db *db.Queries) *userRoleRepository {
	return &userRoleRepository{
		db: db,
	}
}

// AssignRoleToUser assigns a role to a user.
//
// Parameters:
//   - ctx: the context for the database operation
//   - req: the request payload containing user ID and role ID
//
// Returns:
//   - The created UserRoleRecord if successful, or an error if the operation fails.
func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error) {
	res, err := r.db.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}


	return res, nil
}

// RemoveRoleFromUser removes a role assigned to a user.
//
// Parameters:
//   - ctx: the context for the database operation
//   - req: the request payload containing user ID and role ID
//
// Returns:
//   - An error if the operation fails.
func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	err := r.db.RemoveRoleFromUser(ctx, db.RemoveRoleFromUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}


	return nil
}
