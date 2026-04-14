package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
)


// userRepository is a struct that represents a user repository
type userRepository struct {
	db *db.Queries
}

// NewUserRepository returns a new instance of userRepository.
//
// It takes in a db.Queries instance as its database handler, a context.Context
// for its database operations, and a recordmapper.UserRecordMapping for mapper
// database records to the domain level record.UserRecord.
//
// It returns a new instance of userRepository for use in the application.
func NewUserRepository(db *db.Queries) *userRepository {
	return &userRepository{
		db: db,
	}
}

// FindById retrieves a user by their unique ID.
//
// Parameters:
//   - ctx: the context for the database operation
//   - id: the user's unique identifier
//
// Returns:
//   - A UserRecord if found, or an error if the operation fails or user is not found.
func (r *userRepository) FindById(ctx context.Context, user_id int) (*db.GetUserByIDRow, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}


	return res, nil
}


// FindByEmail retrieves a user by their email address.
//
// Parameters:
//   - ctx: the context for the database operation
//   - email: the email address to search for
//
// Returns:
//   - A UserRecord if found, or an error if the operation fails or user does not exist.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*db.User, error) {
	res, err := r.db.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}


	return res, nil
}


// FindByEmailAndVerify retrieves a verified user by their email address.
//
// Parameters:
//   - ctx: the context for the database operation
//   - email: the email address to search for
//
// Returns:
//   - A UserRecord if found and verified, or an error otherwise.
func (r *userRepository) FindByEmailAndVerify(ctx context.Context, email string) (*db.GetUserByEmailAndVerifyRow, error) {
	res, err := r.db.GetUserByEmailAndVerify(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}


	return res, nil
}


// FindByVerificationCode retrieves a user by their verification code.
//
// Parameters:
//   - ctx: the context for the database operation
//   - verificationCode: the verification code string
//
// Returns:
//   - A UserRecord if found, or an error if invalid or not found.
func (r *userRepository) FindByVerificationCode(ctx context.Context, verification_code string) (*db.GetUserByVerificationCodeRow, error) {
	res, err := r.db.GetUserByVerificationCode(ctx, verification_code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}


	return res, nil
}


// CreateUser inserts a new user into the database.
//
// Parameters:
//   - ctx: the context for the database operation
//   - request: the user registration data
//
// Returns:
//   - The created UserRecord, or an error if the operation fails.
func (r *userRepository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*db.CreateUserRow, error) {
	isVerified := true

	req := db.CreateUserParams{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       &isVerified,
	}

	user, err := r.db.CreateUser(ctx, req)

	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}

	return user, nil
}


// UpdateUserIsVerified updates the verification status of a user.
//
// Parameters:
//   - ctx: the context for the database operation
//   - userID: the user's ID
//   - isVerified: the updated verification status
//
// Returns:
//   - The updated UserRecord, or an error if the operation fails.
func (r *userRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*db.UpdateUserIsVerifiedRow, error) {
	res, err := r.db.UpdateUserIsVerified(ctx, db.UpdateUserIsVerifiedParams{
		UserID:     int32(user_id),
		IsVerified: &is_verified,
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}

	return res, nil
}


// UpdateUserPassword updates a user's password.
//
// Parameters:
//   - ctx: the context for the database operation
//   - userID: the user's ID
//   - password: the new password (hashed)
//
// Returns:
//   - The updated UserRecord, or an error if the operation fails.
func (r *userRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*db.UpdateUserPasswordRow, error) {
	res, err := r.db.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		UserID:   int32(user_id),
		Password: password,
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}

	return res, nil
}

