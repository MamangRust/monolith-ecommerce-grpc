package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	userrole_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/user_role_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

// userRoleRepository is a struct that implements the UserRoleRepository interface using gRPC client
type userRoleRepository struct {
	client pb.RoleCommandServiceClient
}

// NewUserRoleRepository creates a new UserRoleRepository instance
func NewUserRoleRepository(client pb.RoleCommandServiceClient) UserRoleRepository {
	return &userRoleRepository{
		client: client,
	}
}

// AssignRoleToUser assigns a role to a user via gRPC.
func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error) {
	protoReq := &pb.AssignRoleToUserRequest{
		UserId: int32(req.UserId),
		RoleId: int32(req.RoleId),
	}

	res, err := r.client.AssignRoleToUser(ctx, protoReq)
	if err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}

	return &db.UserRole{
		UserRoleID: res.Data.UserRoleId,
		UserID:     res.Data.UserId,
		RoleID:     res.Data.RoleId,
	}, nil
}

// RemoveRoleFromUser removes a role assigned to a user via gRPC.
func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	protoReq := &pb.RemoveRoleFromUserRequest{
		UserId: int32(req.UserId),
		RoleId: int32(req.RoleId),
	}

	_, err := r.client.RemoveRoleFromUser(ctx, protoReq)
	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}

	return nil
}
