// Package role adapts the Role service gRPC query API into the shared
// pkg/database sqlc row types.
package role

import (
	"context"
	"time"

	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	role_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"github.com/jackc/pgx/v5/pgtype"
)

// QueryRepository is the contract consumers depend on for role reads.
type QueryRepository interface {
	FindByID(ctx context.Context, roleID int) (*db.Role, error)
	FindById(ctx context.Context, roleID int) (*db.Role, error)
	FindByName(ctx context.Context, name string) (*db.Role, error)
}

// Repository implements QueryRepository over the generated role query client.
type Repository struct {
	client pbrole.RoleQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a role adapter. Passing zero options leaves the guard nil, which
// makes DependencyGuard.Call a plain passthrough.
func New(client pbrole.RoleQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewAdapter is an alias for New, matching the naming used by the other adapters.
func NewAdapter(client pbrole.RoleQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(client, opts...)
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindByID(ctx context.Context, roleID int) (*db.Role, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbrole.ApiResponseRole, error) {
		return r.client.FindByIdRole(callCtx, &pbrole.FindByIdRoleRequest{RoleId: int32(roleID)})
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, role_errors.ErrRoleNotFound
	}
	return toDBRole(resp.Data), nil
}

func (r *Repository) FindById(ctx context.Context, roleID int) (*db.Role, error) {
	return r.FindByID(ctx, roleID)
}

func (r *Repository) FindByName(ctx context.Context, name string) (*db.Role, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbrole.ApiResponseRole, error) {
		return r.client.FindByNameRole(callCtx, &pbrole.FindByNameRoleRequest{Name: name})
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, role_errors.ErrRoleNotFound
	}
	return toDBRole(resp.Data), nil
}

func toDBRole(role *pbrole.RoleResponse) *db.Role {
	if role == nil {
		return nil
	}

	createdAt, cErr := time.Parse(time.RFC3339, role.CreatedAt)
	updatedAt, uErr := time.Parse(time.RFC3339, role.UpdatedAt)

	return &db.Role{
		RoleID:    role.Id,
		RoleName:  role.Name,
		CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: cErr == nil},
		UpdatedAt: pgtype.Timestamp{Time: updatedAt, Valid: uErr == nil},
	}
}
