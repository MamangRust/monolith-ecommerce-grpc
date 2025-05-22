package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type roleRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.RoleRecordMapping
}

func NewRoleRepository(db *db.Queries, ctx context.Context, mapping recordmapper.RoleRecordMapping) *roleRepository {
	return &roleRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *roleRepository) FindById(id int) (*record.RoleRecord, error) {
	res, err := r.db.GetRole(r.ctx, int32(id))
	if err != nil {
		return nil, role_errors.ErrRoleNotFound
	}
	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleRepository) FindByName(name string) (*record.RoleRecord, error) {
	res, err := r.db.GetRoleByName(r.ctx, name)
	if err != nil {
		return nil, role_errors.ErrRoleNotFound
	}
	return r.mapping.ToRoleRecord(res), nil
}
