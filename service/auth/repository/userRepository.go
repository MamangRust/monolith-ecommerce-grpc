package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// userRepository is a struct that implements the UserRepository interface using gRPC clients
type userRepository struct {
	queryClient   pb.UserQueryServiceClient
	commandClient pb.UserCommandServiceClient
}

// NewUserRepository returns a new instance of userRepository using gRPC clients.
func NewUserRepository(queryClient pb.UserQueryServiceClient, commandClient pb.UserCommandServiceClient) *userRepository {
	return &userRepository{
		queryClient:   queryClient,
		commandClient: commandClient,
	}
}

// FindById retrieves a user by their unique ID via gRPC.
func (r *userRepository) FindById(ctx context.Context, user_id int) (*db.GetUserByIDRow, error) {
	res, err := r.queryClient.FindById(ctx, &pb.FindByIdUserRequest{Id: int32(user_id)})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &db.GetUserByIDRow{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}

// FindByEmail retrieves a user by their email address via gRPC.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*db.User, error) {
	res, err := r.queryClient.FindByEmail(ctx, &pb.FindByEmailRequest{Email: email})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, nil
		}
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &db.User{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
		Password:  res.Data.Password,
	}, nil
}

// FindByEmailAndVerify retrieves a verified user by their email address via gRPC.
func (r *userRepository) FindByEmailAndVerify(ctx context.Context, email string) (*db.GetUserByEmailAndVerifyRow, error) {
	// We use FindByEmail but we should probably check if it's verified in the auth service or add a Verify check in the request.
	// For now, our FindByEmail returns the user if found.
	res, err := r.queryClient.FindByEmail(ctx, &pb.FindByEmailRequest{Email: email})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &db.GetUserByEmailAndVerifyRow{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
		Password:  res.Data.Password,
	}, nil
}

// FindByVerificationCode retrieves a user by their verification code via gRPC.
func (r *userRepository) FindByVerificationCode(ctx context.Context, verification_code string) (*db.GetUserByVerificationCodeRow, error) {
	res, err := r.queryClient.FindByVerificationCode(ctx, &pb.FindByVerificationCodeRequest{VerificationCode: verification_code})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &db.GetUserByVerificationCodeRow{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}

// CreateUser inserts a new user via gRPC.
func (r *userRepository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*db.CreateUserRow, error) {
	res, err := r.commandClient.Create(ctx, &pb.CreateUserRequest{
		Firstname:       request.FirstName,
		Lastname:        request.LastName,
		Email:           request.Email,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	})

	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}

	return &db.CreateUserRow{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}

// UpdateUserIsVerified updates the verification status of a user via gRPC.
func (r *userRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*db.UpdateUserIsVerifiedRow, error) {
	res, err := r.commandClient.UpdateIsVerified(ctx, &pb.UpdateUserIsVerifiedRequest{
		Id:         int32(user_id),
		IsVerified: is_verified,
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}

	return &db.UpdateUserIsVerifiedRow{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}

// UpdateUserPassword updates a user's password via gRPC.
func (r *userRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*db.UpdateUserPasswordRow, error) {
	res, err := r.commandClient.UpdatePassword(ctx, &pb.UpdateUserPasswordRequest{
		Id:       int32(user_id),
		Password: password,
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}

	return &db.UpdateUserPasswordRow{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}
