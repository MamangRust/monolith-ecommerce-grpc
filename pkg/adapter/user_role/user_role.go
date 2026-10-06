// Package user_role adapts the UserRole service gRPC command API into the
// shared pkg/database sqlc row types.
package user_role

import (
	"context"

	pbuserrole "github.com/MamangRust/monolith-ecommerce-pb/user_role"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	userrole_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/user_role_errors"
)

// CommandRepository is the write path consumers use to (un)assign roles.
type CommandRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// Repository implements CommandRepository on top of the generated user-role
// service client.
type Repository struct {
	client pbuserrole.UserRoleServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a user-role adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough.
func New(client pbuserrole.UserRoleServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

// AssignRoleToUser implements CommandRepository.
func (r *Repository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuserrole.ApiResponseUserRole, error) {
		return r.client.AssignRoleToUser(callCtx, &pbuserrole.AssignRoleToUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
	})
	if err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, userrole_errors.ErrAssignRoleToUser
	}

	return &db.UserRole{
		UserRoleID: resp.Data.UserRoleId,
		UserID:     resp.Data.UserId,
		RoleID:     resp.Data.RoleId,
	}, nil
}

// RemoveRoleFromUser implements CommandRepository.
func (r *Repository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	_, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (struct{}, error) {
		_, callErr := r.client.RemoveRoleFromUser(callCtx, &pbuserrole.RemoveRoleFromUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
		return struct{}{}, callErr
	})
	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}
	return nil
}
