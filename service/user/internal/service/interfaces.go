package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type UserQueryService interface {
	FindAll(req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse)
	FindByID(id int) (*response.UserResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)
}

type UserCommandService interface {
	CreateUser(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	UpdateUser(request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	TrashedUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	RestoreUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	DeleteUserPermanent(user_id int) (bool, *response.ErrorResponse)

	RestoreAllUser() (bool, *response.ErrorResponse)
	DeleteAllUserPermanent() (bool, *response.ErrorResponse)
}
