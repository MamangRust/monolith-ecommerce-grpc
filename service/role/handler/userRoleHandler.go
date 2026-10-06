package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-role/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"google.golang.org/protobuf/types/known/emptypb"

	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	pbuserrole "github.com/MamangRust/monolith-ecommerce-pb/user_role"
)

type userRoleHandler struct {
	pbuserrole.UnimplementedUserRoleServiceServer
	roleQuery   service.RoleQueryService
	roleCommand service.RoleCommandService
	logger      logger.LoggerInterface
}

// NewUserRoleHandler exposes the user-role operations as their own gRPC service
// (which piggybacks on the role service). It delegates to the already-existing
// role query/command services so there is a single source of truth.
func NewUserRoleHandler(svc *service.Service, logger logger.LoggerInterface) pbuserrole.UserRoleServiceServer {
	return &userRoleHandler{
		roleQuery:   svc.RoleQuery,
		roleCommand: svc.RoleCommand,
		logger:      logger,
	}
}

func (s *userRoleHandler) FindByUserId(ctx context.Context, req *pbuserrole.FindByIdUserRoleRequest) (*pbrole.ApiResponsesRole, error) {
	userID := int(req.GetUserId())
	if userID == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	roles, err := s.roleQuery.FindByUserId(ctx, userID)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoRoles := make([]*pbrole.RoleResponse, len(roles))
	for i, role := range roles {
		protoRoles[i] = mapToProtoRoleResponse(role)
	}

	return &pbrole.ApiResponsesRole{
		Status:  "success",
		Message: "Successfully fetched role by user id",
		Data:    protoRoles,
	}, nil
}

func (s *userRoleHandler) AssignRoleToUser(ctx context.Context, req *pbuserrole.AssignRoleToUserRequest) (*pbuserrole.ApiResponseUserRole, error) {
	userID := int(req.GetUserId())
	roleID := int(req.GetRoleId())
	if userID == 0 || roleID == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	userRole, err := s.roleCommand.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbuserrole.ApiResponseUserRole{
		Status:  "success",
		Message: "Successfully assigned role to user",
		Data: &pbuserrole.UserRoleResponse{
			UserRoleId: userRole.UserRoleID,
			UserId:     userRole.UserID,
			RoleId:     userRole.RoleID,
		},
	}, nil
}

func (s *userRoleHandler) RemoveRoleFromUser(ctx context.Context, req *pbuserrole.RemoveRoleFromUserRequest) (*emptypb.Empty, error) {
	userID := int(req.GetUserId())
	roleID := int(req.GetRoleId())
	if userID == 0 || roleID == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	if err := s.roleCommand.RemoveRoleFromUser(ctx, &requests.RemoveUserRoleRequest{
		UserId: userID,
		RoleId: roleID,
	}); err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}
