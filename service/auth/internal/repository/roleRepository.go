package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
)


// roleRepository is a struct that implements the RoleRepository interface
type roleRepository struct {
	db *db.Queries
}

// NewRoleRepository creates a new RoleRepository instance
//
// Args:
// db: a pointer to the database queries
// ctx: a context.Context object
// mapper: a RoleRecordMapping object
//
// Returns:
// a pointer to the roleRepository struct
func NewRoleRepository(db *db.Queries) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

// FindById retrieves a role by its unique ID.
//
// Parameters:
//   - ctx: the context for the database operation
//   - role_id: the unique identifier of the role
//
// Returns:
//   - A RoleRecord if found, or an error if the role does not exist or operation fails.
func (r *roleRepository) FindById(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.GetRole(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}


// FindByName retrieves a role by its name from the database.
//
// Args:
// name: The name of the role to retrieve.
// // FindByName retrieves a role by its name.
//
// Parameters:
//   - ctx: the context for the database operation
//   - name: the name of the role to search for
//
// Returns:
//   - A RoleRecord if found, or an error if the rol
func (r *roleRepository) FindByName(ctx context.Context, name string) (*db.Role, error) {
	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

