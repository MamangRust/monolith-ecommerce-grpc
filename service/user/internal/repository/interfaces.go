package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindAllUsers(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindByActive(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindByTrashed(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindById(user_id int) (*record.UserRecord, error)
	FindByEmail(email string) (*record.UserRecord, error)
}

type UserCommandRepository interface {
	CreateUser(request *requests.CreateUserRequest) (*record.UserRecord, error)
	UpdateUser(request *requests.UpdateUserRequest) (*record.UserRecord, error)
	TrashedUser(user_id int) (*record.UserRecord, error)
	RestoreUser(user_id int) (*record.UserRecord, error)
	DeleteUserPermanent(user_id int) (bool, error)
	RestoreAllUser() (bool, error)
	DeleteAllUserPermanent() (bool, error)
}

type RoleRepository interface {
	FindById(role_id int) (*record.RoleRecord, error)
	FindByName(name string) (*record.RoleRecord, error)
}
