package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userHandleGrpc struct {
	pb.UnimplementedUserServiceServer
	userQueryService   service.UserQueryService
	userCommandService service.UserCommandService
	mapping            protomapper.UserProtoMapper
}

func NewUserHandleGrpc(user service.Service) *userHandleGrpc {
	return &userHandleGrpc{
		userQueryService:   user.UserQuery,
		userCommandService: user.UserCommand,
		mapping:            protomapper.NewUserProtoMapper(),
	}
}

func (s *userHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUser, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationUser(paginationMeta, "success", "Successfully fetched users", users)
	return so, nil
}

func (s *userHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, user_errors.ErrGrpcUserNotFound
	}

	user, err := s.userQueryService.FindByID(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUser("success", "Successfully fetched user", user)

	return so, nil

}

func (s *userHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByActive(&reqService)

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
	so := s.mapping.ToProtoResponsePaginationUserDeleteAt(paginationMeta, "success", "Successfully fetched active users", users)

	return so, nil
}

func (s *userHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationUserDeleteAt(paginationMeta, "success", "Successfully fetched trashed users", users)

	return so, nil
}

func (s *userHandleGrpc) Create(ctx context.Context, request *pb.CreateUserRequest) (*pb.ApiResponseUser, error) {
	req := &requests.CreateUserRequest{
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		return nil, user_errors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.CreateUser(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUser("success", "Successfully created user", user)

	return so, nil
}

func (s *userHandleGrpc) Update(ctx context.Context, request *pb.UpdateUserRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	req := &requests.UpdateUserRequest{
		UserID:          &id,
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		return nil, user_errors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.UpdateUser(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUser("success", "Successfully updated user", user)

	return so, nil
}

func (s *userHandleGrpc) TrashedUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.TrashedUser(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUserDeleteAt("success", "Successfully trashed user", user)

	return so, nil
}

func (s *userHandleGrpc) RestoreUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.RestoreUser(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUserDeleteAt("success", "Successfully restored user", user)

	return so, nil
}

func (s *userHandleGrpc) DeleteUserPermanent(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	_, err := s.userCommandService.DeleteUserPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUserDelete("success", "Successfully deleted user permanently")

	return so, nil
}

func (s *userHandleGrpc) RestoreAllUser(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	_, err := s.userCommandService.RestoreAllUser()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUserAll("success", "Successfully restore all user")

	return so, nil
}

func (s *userHandleGrpc) DeleteAllUserPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	_, err := s.userCommandService.DeleteAllUserPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUserAll("success", "Successfully delete user permanen")

	return so, nil
}
