package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

// UserQueryService handles query operations related to user data.
type UserQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, *int, error)
	FindByID(ctx context.Context, id int) (*db.GetUserByIDRow, error)
	FindByEmailWithPassword(ctx context.Context, email string) (*db.GetUserByEmailWithPasswordRow, error)
	FindByVerificationCode(ctx context.Context, code string) (*db.GetUserByVerificationCodeRow, error)
	FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, *int, error)
}

// UserCommandService handles command operations related to user management.
type UserCommandService interface {
	Create(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error)
	Update(ctx context.Context, request *requests.UpdateUserRequest) (*db.User, error)
	UpdateIsVerified(ctx context.Context, user_id int, is_verified bool) (*db.User, error)
	UpdatePassword(ctx context.Context, user_id int, password string) (*db.User, error)
	Trash(ctx context.Context, user_id int) (*db.TrashUserRow, error)
	Restore(ctx context.Context, user_id int) (*db.RestoreUserRow, error)
	DeletePermanent(ctx context.Context, user_id int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
