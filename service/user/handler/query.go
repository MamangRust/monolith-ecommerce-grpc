package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type userQueryHandler struct {
	pb.UnimplementedUserQueryServiceServer
	UserQuery service.UserQueryService
	logger    logger.LoggerInterface
}

func NewUserQueryHandler(svc service.UserQueryService, logger logger.LoggerInterface) UserQueryHandler {
	return &userQueryHandler{
		UserQuery: svc,
		logger:    logger,
	}
}

func (s *userQueryHandler) FindAll(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUser, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.UserQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoUsers := make([]*pb.UserResponse, len(users))
	for i, user := range users {
		protoUsers[i] = mapToProtoUserResponse(user)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationUser{
		Status:     "success",
		Message:    "Successfully fetched users",
		Data:       protoUsers,
		Pagination: paginationMeta,
	}, nil
}

func (s *userQueryHandler) FindById(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.UserQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully fetched user",
		Data:    mapToProtoUserResponse(user),
	}, nil
}

func (s *userQueryHandler) FindByActive(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.UserQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoUsers := make([]*pb.UserResponseDeleteAt, len(users))
	for i, user := range users {
		protoUsers[i] = mapToProtoUserResponseDeleteAt(user)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationUserDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active users",
		Data:       protoUsers,
		Pagination: paginationMeta,
	}, nil
}

func (s *userQueryHandler) FindByTrashed(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.UserQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoUsers := make([]*pb.UserResponseDeleteAt, len(users))
	for i, user := range users {
		protoUsers[i] = mapToProtoUserResponseDeleteAt(user)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationUserDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed users",
		Data:       protoUsers,
		Pagination: paginationMeta,
	}, nil
}
func (s *userQueryHandler) FindByEmail(ctx context.Context, request *pb.FindByEmailRequest) (*pb.ApiResponseUserWithPassword, error) {
	email := request.GetEmail()
	if email == "" {
		return nil, user_errors.ErrGrpcUserInvalidEmail
	}

	user, err := s.UserQuery.FindByEmailWithPassword(ctx, email)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUserWithPassword{
		Status:  "success",
		Message: "Successfully fetched user by email",
		Data:    mapToProtoUserResponseWithPassword(user),
	}, nil
}

func (s *userQueryHandler) FindByVerificationCode(ctx context.Context, request *pb.FindByVerificationCodeRequest) (*pb.ApiResponseUser, error) {
	code := request.GetVerificationCode()
	if code == "" {
		return nil, user_errors.ErrGrpcUserInvalidVerificationCode
	}

	user, err := s.UserQuery.FindByVerificationCode(ctx, code)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully fetched user by verification code",
		Data:    mapToProtoUserResponse(user),
	}, nil
}
