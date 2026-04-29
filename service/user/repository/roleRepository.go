package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	role_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type roleRepository struct {
	client pb.RoleQueryServiceClient
}

func NewRoleRepository(client pb.RoleQueryServiceClient) *roleRepository {
	return &roleRepository{
		client: client,
	}
}

func (r *roleRepository) FindByID(ctx context.Context, role_id int) (*db.Role, error) {
	res, err := r.client.FindByIdRole(ctx, &pb.FindByIdRoleRequest{RoleId: int32(role_id)})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}

	return &db.Role{
		RoleID:   res.Data.Id,
		RoleName: res.Data.Name,
	}, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*db.Role, error) {
	res, err := r.client.FindAllRole(ctx, &pb.FindAllRoleRequest{
		Search:   name,
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}

	for _, role := range res.Data {
		if role.Name == name {
			return &db.Role{
				RoleID:   role.Id,
				RoleName: role.Name,
			}, nil
		}
	}

	return nil, errors.ErrNotFound.WithMessage("role not found by name: " + name)
}
