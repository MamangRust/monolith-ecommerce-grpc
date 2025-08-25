package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userHandleGrpc struct {
	pb.UnimplementedUserServiceServer
	userQueryService   service.UserQueryService
	userCommandService service.UserCommandService
	logger             logger.LoggerInterface
	mapping            protomapper.UserProtoMapper
}

func NewUserHandleGrpc(user *service.Service, logger logger.LoggerInterface) pb.UserServiceServer {
	return &userHandleGrpc{
		userQueryService:   user.UserQuery,
		userCommandService: user.UserCommand,
		logger:             logger,
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

	s.logger.Info("Fetching all users",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all users",
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

	s.logger.Info("Successfully fetched all users",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_users_count", len(users)),
	)

	so := s.mapping.ToProtoResponsePaginationUser(paginationMeta, "success", "Successfully fetched users", users)
	return so, nil
}

func (s *userHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid user ID provided", zap.Int("user_id", id))
		return nil, user_errors.ErrGrpcUserNotFound
	}

	s.logger.Info("Fetching user by ID", zap.Int("user_id", id))

	user, err := s.userQueryService.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch user by ID",
			zap.Int("user_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched user by ID",
		zap.Int("user_id", id),
	)

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

	s.logger.Info("Fetching active users",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active users",
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

	s.logger.Info("Successfully fetched active users",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_users_count", len(users)),
	)

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

	s.logger.Info("Fetching trashed users",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed users",
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

	s.logger.Info("Successfully fetched trashed users",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_users_count", len(users)),
	)

	so := s.mapping.ToProtoResponsePaginationUserDeleteAt(paginationMeta, "success", "Successfully fetched trashed users", users)
	return so, nil
}

func (s *userHandleGrpc) Create(ctx context.Context, request *pb.CreateUserRequest) (*pb.ApiResponseUser, error) {
	s.logger.Info("Creating new user",
		zap.String("email", request.GetEmail()),
		zap.String("first_name", request.GetFirstname()),
		zap.String("last_name", request.GetLastname()),
	)

	req := &requests.CreateUserRequest{
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on user creation",
			zap.String("email", request.GetEmail()),
			zap.Error(err),
		)
		return nil, user_errors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.CreateUser(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create user",
			zap.String("email", request.GetEmail()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("User created successfully",
		zap.Int("user_id", int(user.ID)),
		zap.String("email", user.Email),
		zap.String("full_name", user.FirstName+" "+user.LastName),
	)

	so := s.mapping.ToProtoResponseUser("success", "Successfully created user", user)
	return so, nil
}

func (s *userHandleGrpc) Update(ctx context.Context, request *pb.UpdateUserRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid user ID provided for update", zap.Int("user_id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	s.logger.Info("Updating user", zap.Int("user_id", id))

	req := &requests.UpdateUserRequest{
		UserID:          &id,
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on user update",
			zap.Int("user_id", id),
			zap.String("email", request.GetEmail()),
			zap.Error(err),
		)
		return nil, user_errors.ErrGrpcValidateUpdateUser
	}

	user, err := s.userCommandService.UpdateUser(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update user",
			zap.Int("user_id", id),
			zap.String("email", request.GetEmail()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("User updated successfully",
		zap.Int("user_id", id),
		zap.String("full_name", user.FirstName+" "+user.LastName),
		zap.String("email", user.Email),
	)

	so := s.mapping.ToProtoResponseUser("success", "Successfully updated user", user)
	return so, nil
}

func (s *userHandleGrpc) TrashedUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid user ID for trashing", zap.Int("user_id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	s.logger.Info("Moving user to trash", zap.Int("user_id", id))

	user, err := s.userCommandService.TrashedUser(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash user",
			zap.Int("user_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("User moved to trash successfully",
		zap.Int("user_id", id),
		zap.String("email", user.Email),
		zap.String("full_name", user.FirstName+" "+user.LastName),
	)

	so := s.mapping.ToProtoResponseUserDeleteAt("success", "Successfully trashed user", user)
	return so, nil
}

func (s *userHandleGrpc) RestoreUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid user ID for restore", zap.Int("user_id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	s.logger.Info("Restoring user from trash", zap.Int("user_id", id))

	user, err := s.userCommandService.RestoreUser(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore user",
			zap.Int("user_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("User restored successfully",
		zap.Int("user_id", id),
		zap.String("email", user.Email),
		zap.String("full_name", user.FirstName+" "+user.LastName),
	)

	so := s.mapping.ToProtoResponseUserDeleteAt("success", "Successfully restored user", user)
	return so, nil
}

func (s *userHandleGrpc) DeleteUserPermanent(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid user ID for permanent deletion", zap.Int("user_id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	s.logger.Info("Permanently deleting user", zap.Int("user_id", id))

	_, err := s.userCommandService.DeleteUserPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete user",
			zap.Int("user_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("User permanently deleted", zap.Int("user_id", id))

	so := s.mapping.ToProtoResponseUserDelete("success", "Successfully deleted user permanently")
	return so, nil
}

func (s *userHandleGrpc) RestoreAllUser(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	s.logger.Info("Restoring all trashed users")

	_, err := s.userCommandService.RestoreAllUser(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all users", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All users restored successfully")

	so := s.mapping.ToProtoResponseUserAll("success", "Successfully restored all users")
	return so, nil
}

func (s *userHandleGrpc) DeleteAllUserPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	s.logger.Info("Permanently deleting all trashed users")

	_, err := s.userCommandService.DeleteAllUserPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all users", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All users permanently deleted")

	so := s.mapping.ToProtoResponseUserAll("success", "Successfully deleted all users permanently")
	return so, nil
}
