package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, error)
	FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, error)
	FindByID(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
	FindByEmail(ctx context.Context, email string) (*db.User, error)
	FindByIDWithPassword(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
	FindByEmailWithPassword(ctx context.Context, email string) (*db.GetUserByEmailWithPasswordRow, error)
	FindByVerificationCode(ctx context.Context, code string) (*db.GetUserByVerificationCodeRow, error)
}

type UserCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error)
	Update(ctx context.Context, request *requests.UpdateUserRequest) (*db.User, error)
	UpdateIsVerified(ctx context.Context, user_id int, is_verified bool) (*db.UpdateUserIsVerifiedRow, error)
	UpdatePassword(ctx context.Context, user_id int, password string) (*db.UpdateUserPasswordRow, error)
	Trash(ctx context.Context, user_id int) (*db.TrashUserRow, error)
	Restore(ctx context.Context, user_id int) (*db.RestoreUserRow, error)
	DeletePermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type RoleRepository interface {
	FindByID(ctx context.Context, role_id int) (*db.Role, error)
	FindByName(ctx context.Context, name string) (*db.Role, error)
}
