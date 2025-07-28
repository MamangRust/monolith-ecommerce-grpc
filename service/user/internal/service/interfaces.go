package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type UserQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse)
	FindByID(ctx context.Context, id int) (*response.UserResponse, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)
}

type UserCommandService interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	TrashedUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	RestoreUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, *response.ErrorResponse)

	RestoreAllUser(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllUserPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
