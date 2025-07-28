package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type roleRepository struct {
	db      *db.Queries
	mapping recordmapper.RoleRecordMapping
}

func NewRoleRepository(db *db.Queries, mapping recordmapper.RoleRecordMapping) *roleRepository {
	return &roleRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *roleRepository) FindById(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.GetRole(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found with ID: %d", id)
		}
		return nil, fmt.Errorf("failed to find role by ID %d: %w", id, err)
	}
	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*record.RoleRecord, error) {
	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound
		}

		return nil, role_errors.ErrRoleNotFound
	}
	return r.mapping.ToRoleRecord(res), nil
}
