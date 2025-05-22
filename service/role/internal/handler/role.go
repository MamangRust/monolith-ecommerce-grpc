package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleHandleGrpc struct {
	pb.UnimplementedRoleServiceServer
	roleQuery   service.RoleQueryService
	roleCommand service.RoleCommandService
	mapping     protomapper.RoleProtoMapper
}

func NewRoleHandleGrpc(service service.Service) *roleHandleGrpc {
	return &roleHandleGrpc{
		roleQuery:   service.RoleQuery,
		roleCommand: service.RoleCommand,
		mapping:     protomapper.NewRoleProtoMapper(),
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

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	role, totalRecords, err := s.roleQuery.FindAll(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationRole(paginationMeta, "success", "Successfully fetched role records", role)

	return so, nil
}

func (s *roleHandleGrpc) FindByIdRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	if roleID == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleQuery.FindById(roleID)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	roleResponse := s.mapping.ToProtoResponseRole("success", "Successfully fetched role", role)

	return roleResponse, nil
}

func (s *roleHandleGrpc) FindByUserId(ctx context.Context, req *pb.FindByIdUserRoleRequest) (*pb.ApiResponsesRole, error) {
	userID := int(req.GetUserId())

	if userID == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleQuery.FindByUserId(userID)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	roleResponse := s.mapping.ToProtoResponsesRole("success", "Successfully fetched role by user ID", role)

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

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByActiveRole(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
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

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByTrashedRole(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched trashed roles", roles)

	return so, nil
}

func (s *roleHandleGrpc) CreateRole(ctx context.Context, reqPb *pb.CreateRoleRequest) (*pb.ApiResponseRole, error) {
	req := &requests.CreateRoleRequest{
		Name: reqPb.Name,
	}

	if err := req.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateCreateRole
	}

	role, err := s.roleCommand.CreateRole(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully created role", role)

	return so, nil
}

func (s *roleHandleGrpc) UpdateRole(ctx context.Context, reqPb *pb.UpdateRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(reqPb.GetId())

	if roleID == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	name := reqPb.GetName()

	req := &requests.UpdateRoleRequest{
		ID:   &roleID,
		Name: name,
	}

	if err := req.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateUpdateRole
	}

	role, err := s.roleCommand.UpdateRole(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully updated role", role)

	return so, nil
}

func (s *roleHandleGrpc) TrashedRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	if roleID == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.TrashedRole(roleID)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully trashed role", role)

	return so, nil
}

func (s *roleHandleGrpc) RestoreRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	if roleID == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.RestoreRole(roleID)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully restored role", role)

	return so, nil
}

func (s *roleHandleGrpc) DeleteRolePermanent(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRoleDelete, error) {
	id := int(req.GetRoleId())

	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	_, err := s.roleCommand.DeleteRolePermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRoleDelete("success", "Successfully deleted role permanently")

	return so, nil
}

func (s *roleHandleGrpc) RestoreAllRole(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	_, err := s.roleCommand.RestoreAllRole()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRoleAll("success", "Successfully restored all roles")

	return so, nil
}

func (s *roleHandleGrpc) DeleteAllRolePermanent(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	_, err := s.roleCommand.DeleteAllRolePermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRoleAll("success", "Successfully deleted all roles")

	return so, nil
}
