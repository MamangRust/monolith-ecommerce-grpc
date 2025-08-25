package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleHandleGrpc struct {
	pb.UnimplementedRoleServiceServer
	roleQuery   service.RoleQueryService
	roleCommand service.RoleCommandService
	mapping     protomapper.RoleProtoMapper
	logger      logger.LoggerInterface
}

func NewRoleHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.RoleServiceServer {
	return &roleHandleGrpc{
		roleQuery:   service.RoleQuery,
		roleCommand: service.RoleCommand,
		mapping:     protomapper.NewRoleProtoMapper(),
		logger:      logger,
	}
}

func (s *roleHandleGrpc) FindAllRole(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRole, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all roles",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all roles",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched all roles",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_roles_count", len(roles)),
	)

	so := s.mapping.ToProtoResponsePaginationRole(paginationMeta, "success", "Successfully fetched role records", roles)
	return so, nil
}

func (s *roleHandleGrpc) FindByIdRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	if roleID == 0 {
		s.logger.Error("Invalid role ID provided", zap.Int("role_id", roleID))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	s.logger.Info("Fetching role by ID", zap.Int("role_id", roleID))

	role, err := s.roleQuery.FindById(ctx, roleID)
	if err != nil {
		s.logger.Error("Failed to fetch role by ID",
			zap.Int("role_id", roleID),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched role by ID",
		zap.Int("role_id", roleID),
		zap.String("role_name", role.Name),
	)

	roleResponse := s.mapping.ToProtoResponseRole("success", "Successfully fetched role", role)
	return roleResponse, nil
}

func (s *roleHandleGrpc) FindByUserId(ctx context.Context, req *pb.FindByIdUserRoleRequest) (*pb.ApiResponsesRole, error) {
	userID := int(req.GetUserId())

	if userID == 0 {
		s.logger.Error("Invalid user ID provided for fetching roles", zap.Int("user_id", userID))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	s.logger.Info("Fetching roles by user ID", zap.Int("user_id", userID))

	roles, err := s.roleQuery.FindByUserId(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to fetch roles by user ID",
			zap.Int("user_id", userID),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched roles by user ID",
		zap.Int("user_id", userID),
		zap.Int("roles_count", len(roles)),
	)

	roleResponse := s.mapping.ToProtoResponsesRole("success", "Successfully fetched role by user ID", roles)
	return roleResponse, nil
}

func (s *roleHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active roles",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByActiveRole(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active roles",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched active roles",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_roles_count", len(roles)),
	)

	so := s.mapping.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched active roles", roles)
	return so, nil
}

func (s *roleHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed roles",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByTrashedRole(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed roles",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched trashed roles",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_roles_count", len(roles)),
	)

	so := s.mapping.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched trashed roles", roles)
	return so, nil
}

func (s *roleHandleGrpc) CreateRole(ctx context.Context, reqPb *pb.CreateRoleRequest) (*pb.ApiResponseRole, error) {
	s.logger.Info("Creating new role",
		zap.String("role_name", reqPb.Name),
	)

	req := &requests.CreateRoleRequest{
		Name: reqPb.Name,
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on role creation",
			zap.String("role_name", reqPb.Name),
			zap.Error(err),
		)
		return nil, role_errors.ErrGrpcValidateCreateRole
	}

	role, err := s.roleCommand.CreateRole(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create role",
			zap.String("role_name", reqPb.Name),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Role created successfully",
		zap.Int("role_id", int(role.ID)),
		zap.String("role_name", role.Name),
	)

	so := s.mapping.ToProtoResponseRole("success", "Successfully created role", role)
	return so, nil
}

func (s *roleHandleGrpc) UpdateRole(ctx context.Context, reqPb *pb.UpdateRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(reqPb.GetId())

	if roleID == 0 {
		s.logger.Error("Invalid role ID provided for update", zap.Int("role_id", roleID))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	s.logger.Info("Updating role", zap.Int("role_id", roleID))

	req := &requests.UpdateRoleRequest{
		ID:   &roleID,
		Name: reqPb.GetName(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on role update",
			zap.Int("role_id", roleID),
			zap.String("new_name", reqPb.GetName()),
			zap.Error(err),
		)
		return nil, role_errors.ErrGrpcValidateUpdateRole
	}

	role, err := s.roleCommand.UpdateRole(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update role",
			zap.Int("role_id", roleID),
			zap.String("new_name", reqPb.GetName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Role updated successfully",
		zap.Int("role_id", roleID),
		zap.String("role_name", role.Name),
	)

	so := s.mapping.ToProtoResponseRole("success", "Successfully updated role", role)
	return so, nil
}

func (s *roleHandleGrpc) TrashedRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	if roleID == 0 {
		s.logger.Error("Invalid role ID for trashing", zap.Int("role_id", roleID))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	s.logger.Info("Moving role to trash", zap.Int("role_id", roleID))

	role, err := s.roleCommand.TrashedRole(ctx, roleID)
	if err != nil {
		s.logger.Error("Failed to trash role",
			zap.Int("role_id", roleID),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Role moved to trash successfully",
		zap.Int("role_id", roleID),
		zap.String("role_name", role.Name),
	)

	so := s.mapping.ToProtoResponseRole("success", "Successfully trashed role", role)
	return so, nil
}

func (s *roleHandleGrpc) RestoreRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	if roleID == 0 {
		s.logger.Error("Invalid role ID for restore", zap.Int("role_id", roleID))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	s.logger.Info("Restoring role from trash", zap.Int("role_id", roleID))

	role, err := s.roleCommand.RestoreRole(ctx, roleID)
	if err != nil {
		s.logger.Error("Failed to restore role",
			zap.Int("role_id", roleID),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Role restored successfully",
		zap.Int("role_id", roleID),
		zap.String("role_name", role.Name),
	)

	so := s.mapping.ToProtoResponseRole("success", "Successfully restored role", role)
	return so, nil
}

func (s *roleHandleGrpc) DeleteRolePermanent(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRoleDelete, error) {
	id := int(req.GetRoleId())

	if id == 0 {
		s.logger.Error("Invalid role ID for permanent deletion", zap.Int("role_id", id))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	s.logger.Info("Permanently deleting role", zap.Int("role_id", id))

	_, err := s.roleCommand.DeleteRolePermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete role",
			zap.Int("role_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Role permanently deleted", zap.Int("role_id", id))

	so := s.mapping.ToProtoResponseRoleDelete("success", "Successfully deleted role permanently")
	return so, nil
}

func (s *roleHandleGrpc) RestoreAllRole(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	s.logger.Info("Restoring all trashed roles")

	_, err := s.roleCommand.RestoreAllRole(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all roles", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All roles restored successfully")

	so := s.mapping.ToProtoResponseRoleAll("success", "Successfully restored all roles")
	return so, nil
}

func (s *roleHandleGrpc) DeleteAllRolePermanent(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	s.logger.Info("Permanently deleting all trashed roles")

	_, err := s.roleCommand.DeleteAllRolePermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all roles", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All roles permanently deleted")

	so := s.mapping.ToProtoResponseRoleAll("success", "Successfully deleted all roles")
	return so, nil
}
