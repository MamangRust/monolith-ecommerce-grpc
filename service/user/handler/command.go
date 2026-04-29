package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userCommandHandler struct {
	pb.UnimplementedUserCommandServiceServer
	UserCommand service.UserCommandService
	logger      logger.LoggerInterface
}

func NewUserCommandHandler(svc service.UserCommandService, logger logger.LoggerInterface) UserCommandHandler {
	return &userCommandHandler{
		UserCommand: svc,
		logger:      logger,
	}
}

func (s *userCommandHandler) Create(ctx context.Context, request *pb.CreateUserRequest) (*pb.ApiResponseUser, error) {
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

	user, err := s.UserCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully created user",
		Data:    mapToProtoUserResponse(user),
	}, nil
}

func (s *userCommandHandler) Update(ctx context.Context, request *pb.UpdateUserRequest) (*pb.ApiResponseUser, error) {
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
		return nil, user_errors.ErrGrpcValidateUpdateUser
	}

	user, err := s.UserCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user",
		Data:    mapToProtoUserResponse(user),
	}, nil
}

func (s *userCommandHandler) TrashedUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.UserCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUserDeleteAt{
		Status:  "success",
		Message: "Successfully trashed user",
		Data:    mapToProtoUserResponseDeleteAt(user),
	}, nil
}

func (s *userCommandHandler) RestoreUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.UserCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUserDeleteAt{
		Status:  "success",
		Message: "Successfully restored user",
		Data:    mapToProtoUserResponseDeleteAt(user),
	}, nil
}

func (s *userCommandHandler) DeleteUserPermanent(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDelete, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	_, err := s.UserCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUserDelete{
		Status:  "success",
		Message: "Successfully deleted user permanently",
	}, nil
}

func (s *userCommandHandler) RestoreAllUser(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	_, err := s.UserCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUserAll{
		Status:  "success",
		Message: "Successfully restored all users",
	}, nil
}

func (s *userCommandHandler) DeleteAllUserPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	_, err := s.UserCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUserAll{
		Status:  "success",
		Message: "Successfully deleted all users permanently",
	}, nil
}
func (s *userCommandHandler) UpdateIsVerified(ctx context.Context, request *pb.UpdateUserIsVerifiedRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.UserCommand.UpdateIsVerified(ctx, id, request.GetIsVerified())
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user verification status",
		Data:    mapToProtoUserResponse(user),
	}, nil
}

func (s *userCommandHandler) UpdatePassword(ctx context.Context, request *pb.UpdateUserPasswordRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.UserCommand.UpdatePassword(ctx, id, request.GetPassword())
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user password",
		Data:    mapToProtoUserResponse(user),
	}, nil
}
