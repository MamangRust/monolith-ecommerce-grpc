package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type UserQueryCache interface {
	GetCachedUsersCache(req *requests.FindAllUsers) ([]*response.UserResponse, *int, bool)
	SetCachedUsersCache(req *requests.FindAllUsers, data []*response.UserResponse, total *int)

	GetCachedUserActiveCache(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool)
	SetCachedUserActiveCache(req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int)

	GetCachedUserTrashedCache(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool)
	SetCachedUserTrashedCache(req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int)

	GetCachedUserCache(id int) (*response.UserResponse, bool)
	SetCachedUserCache(data *response.UserResponse)
}

type UserCommandCache interface {
	DeleteUserCache(id int)
}
