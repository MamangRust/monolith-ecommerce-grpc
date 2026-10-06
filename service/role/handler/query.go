package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-role/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"

	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
)

type roleQueryHandler struct {
	pbrole.UnimplementedRoleQueryServiceServer
	roleQuery service.RoleQueryService
	logger    logger.LoggerInterface
}

func NewRoleQueryHandler(roleQuery service.RoleQueryService, logger logger.LoggerInterface) pbrole.RoleQueryServiceServer {
	return &roleQueryHandler{
		roleQuery: roleQuery,
		logger:    logger,
	}
}

func (s *roleQueryHandler) FindAllRole(ctx context.Context, req *pbrole.FindAllRoleRequest) (*pbrole.ApiResponsePaginationRole, error) {
	page, pageSize := normalizePage(int(req.GetPage()), int(req.GetPageSize()))
	search := req.GetSearch()

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoRoles := make([]*pbrole.RoleResponse, len(roles))
	for i, role := range roles {
		protoRoles[i] = mapToProtoRoleResponse(role)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbrole.ApiResponsePaginationRole{
		Status:     "success",
		Message:    "Successfully fetched role records",
		Data:       protoRoles,
		Pagination: paginationMeta,
	}, nil
}

func (s *roleQueryHandler) FindByActive(ctx context.Context, req *pbrole.FindAllRoleRequest) (*pbrole.ApiResponsePaginationRoleDeleteAt, error) {
	page, pageSize := normalizePage(int(req.GetPage()), int(req.GetPageSize()))
	search := req.GetSearch()

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoRoles := make([]*pbrole.RoleResponseDeleteAt, len(roles))
	for i, role := range roles {
		protoRoles[i] = mapToProtoRoleResponseDeleteAt(role)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbrole.ApiResponsePaginationRoleDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active roles",
		Data:       protoRoles,
		Pagination: paginationMeta,
	}, nil
}

func (s *roleQueryHandler) FindByTrashed(ctx context.Context, req *pbrole.FindAllRoleRequest) (*pbrole.ApiResponsePaginationRoleDeleteAt, error) {
	page, pageSize := normalizePage(int(req.GetPage()), int(req.GetPageSize()))
	search := req.GetSearch()

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoRoles := make([]*pbrole.RoleResponseDeleteAt, len(roles))
	for i, role := range roles {
		protoRoles[i] = mapToProtoRoleResponseDeleteAt(role)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbrole.ApiResponsePaginationRoleDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed roles",
		Data:       protoRoles,
		Pagination: paginationMeta,
	}, nil
}

func (s *roleQueryHandler) FindByIdRole(ctx context.Context, req *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())
	if roleID == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleQuery.FindByID(ctx, roleID)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully fetched role",
		Data:    mapToProtoRoleResponse(role),
	}, nil
}

func (s *roleQueryHandler) FindByNameRole(ctx context.Context, req *pbrole.FindByNameRoleRequest) (*pbrole.ApiResponseRole, error) {
	name := req.GetName()
	if name == "" {
		return nil, role_errors.ErrGrpcRoleInvalidId // Or and appropriate error for empty name
	}

	role, err := s.roleQuery.FindByName(ctx, name)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully fetched role by name",
		Data:    mapToProtoRoleResponse(role),
	}, nil
}

func (s *roleQueryHandler) FindByUserId(ctx context.Context, req *pbrole.FindByIdUserRoleRequest) (*pbrole.ApiResponsesRole, error) {
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
