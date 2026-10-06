// Package user adapts the User service gRPC API into the shared pkg/database
// sqlc row types used by the consumer repositories.
package user

import (
	"context"

	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	user_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueryRepository is the contract consumers depend on for user reads.
type QueryRepository interface {
	FindById(ctx context.Context, userID int) (*db.GetUserByIDRow, error)
	FindByID(ctx context.Context, userID int) (*db.GetUserByIDRow, error)
	FindByEmail(ctx context.Context, email string) (*db.User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*db.GetUserByEmailAndVerifyRow, error)
	FindByVerificationCode(ctx context.Context, code string) (*db.GetUserByVerificationCodeRow, error)
}

// CommandRepository is the contract consumers depend on for user writes.
type CommandRepository interface {
	CreateUser(ctx context.Context, request *requests.RegisterRequest) (*db.CreateUserRow, error)
	UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*db.UpdateUserIsVerifiedRow, error)
	UpdateUserPassword(ctx context.Context, userID int, password string) (*db.UpdateUserPasswordRow, error)
}

// Repository implements QueryRepository and CommandRepository over the generated
// user query/command clients, guarding every call.
type Repository struct {
	query   pbuser.UserQueryServiceClient
	command pbuser.UserCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds a user adapter. Passing zero options leaves the guard nil, which
// makes DependencyGuard.Call a plain passthrough.
func New(query pbuser.UserQueryServiceClient, command pbuser.UserCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter builds a read-only adapter for consumers that never write.
func NewQueryAdapter(query pbuser.UserQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(query, nil, opts...)
}

// NewCommandAdapter builds a write-only adapter for consumers that never read.
func NewCommandAdapter(command pbuser.UserCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return New(nil, command, opts...)
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindById(ctx context.Context, userID int) (*db.GetUserByIDRow, error) {
	return r.findByID(ctx, userID)
}

func (r *Repository) FindByID(ctx context.Context, userID int) (*db.GetUserByIDRow, error) {
	return r.findByID(ctx, userID)
}

func (r *Repository) findByID(ctx context.Context, userID int) (*db.GetUserByIDRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.query.FindById(callCtx, &pbuser.FindByIdUserRequest{Id: int32(userID)})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUserNotFound
	}

	return &db.GetUserByIDRow{
		UserID:    resp.Data.Id,
		Firstname: resp.Data.Firstname,
		Lastname:  resp.Data.Lastname,
		Email:     resp.Data.Email,
	}, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*db.User, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUserWithPassword, error) {
		return r.query.FindByEmail(callCtx, &pbuser.FindByEmailRequest{Email: email})
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, nil
		}
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, nil
	}

	return &db.User{
		UserID:    resp.Data.Id,
		Firstname: resp.Data.Firstname,
		Lastname:  resp.Data.Lastname,
		Email:     resp.Data.Email,
		Password:  resp.Data.Password,
	}, nil
}

func (r *Repository) FindByEmailAndVerify(ctx context.Context, email string) (*db.GetUserByEmailAndVerifyRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUserWithPassword, error) {
		return r.query.FindByEmail(callCtx, &pbuser.FindByEmailRequest{Email: email})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUserNotFound
	}

	return &db.GetUserByEmailAndVerifyRow{
		UserID:    resp.Data.Id,
		Firstname: resp.Data.Firstname,
		Lastname:  resp.Data.Lastname,
		Email:     resp.Data.Email,
		Password:  resp.Data.Password,
	}, nil
}

func (r *Repository) FindByVerificationCode(ctx context.Context, code string) (*db.GetUserByVerificationCodeRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.query.FindByVerificationCode(callCtx, &pbuser.FindByVerificationCodeRequest{VerificationCode: code})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUserNotFound
	}

	return &db.GetUserByVerificationCodeRow{
		UserID:    resp.Data.Id,
		Firstname: resp.Data.Firstname,
		Lastname:  resp.Data.Lastname,
		Email:     resp.Data.Email,
	}, nil
}

func (r *Repository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*db.CreateUserRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.command.Create(callCtx, &pbuser.CreateUserRequest{
			Firstname:       request.FirstName,
			Lastname:        request.LastName,
			Email:           request.Email,
			Password:        request.Password,
			ConfirmPassword: request.ConfirmPassword,
		})
	})
	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrCreateUser
	}

	return &db.CreateUserRow{
		UserID:    resp.Data.Id,
		Firstname: resp.Data.Firstname,
		Lastname:  resp.Data.Lastname,
		Email:     resp.Data.Email,
	}, nil
}

func (r *Repository) UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*db.UpdateUserIsVerifiedRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.command.UpdateIsVerified(callCtx, &pbuser.UpdateUserIsVerifiedRequest{
			Id:         int32(userID),
			IsVerified: isVerified,
		})
	})
	if err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUpdateUserVerificationCode
	}

	return &db.UpdateUserIsVerifiedRow{
		UserID:    resp.Data.Id,
		Firstname: resp.Data.Firstname,
		Lastname:  resp.Data.Lastname,
		Email:     resp.Data.Email,
	}, nil
}

func (r *Repository) UpdateUserPassword(ctx context.Context, userID int, password string) (*db.UpdateUserPasswordRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.command.UpdatePassword(callCtx, &pbuser.UpdateUserPasswordRequest{
			Id:       int32(userID),
			Password: password,
		})
	})
	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUpdateUserPassword
	}

	return &db.UpdateUserPasswordRow{
		UserID:    resp.Data.Id,
		Firstname: resp.Data.Firstname,
		Lastname:  resp.Data.Lastname,
		Email:     resp.Data.Email,
	}, nil
}
