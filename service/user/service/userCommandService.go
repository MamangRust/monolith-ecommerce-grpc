package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type userCommandService struct {
	observability       observability.TraceLoggerObservability
	cache               cache.UserCommandCache
	userCommandRepository repository.UserCommandRepository
	userQueryRepository   repository.UserQueryRepository
	roleRepository      repository.RoleRepository
	logger              logger.LoggerInterface
	hashing             hash.HashPassword
}

type UserCommandServiceDeps struct {
	Observability       observability.TraceLoggerObservability
	Cache               cache.UserCommandCache
	UserCommandRepository repository.UserCommandRepository
	UserQueryRepository   repository.UserQueryRepository
	RoleRepository      repository.RoleRepository
	Logger              logger.LoggerInterface
	Hash                hash.HashPassword
}

func NewUserCommandService(deps *UserCommandServiceDeps) UserCommandService {
	return &userCommandService{
		observability:         deps.Observability,
		cache:                 deps.Cache,
		userCommandRepository: deps.UserCommandRepository,
		userQueryRepository:   deps.UserQueryRepository,
		roleRepository:        deps.RoleRepository,
		logger:                deps.Logger,
		hashing:               deps.Hash,
	}
}

func (s *userCommandService) Create(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Creating new user", zap.String("email", request.Email), zap.Any("request", request))

	existingUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)
	if err == nil && existingUser != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateUserRow](
			s.logger,
			user_errors.ErrUserEmailAlready,
			method,
			span,
			zap.String("email", request.Email),
		)
	}

	hash, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateUserRow](
			s.logger,
			user_errors.ErrUserPassword,
			method,
			span,
		)
	}

	request.Password = hash

	res, err := s.userCommandRepository.Create(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateUserRow](
			s.logger,
			user_errors.ErrFailedCreateUser,
			method,
			span,
		)
	}

	logSuccess("Successfully created new user", zap.String("email", res.Email), zap.Int("user_id", int(res.UserID)))

	return res, nil
}

func (s *userCommandService) Update(ctx context.Context, request *requests.UpdateUserRequest) (*db.User, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", *request.UserID))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Updating user", zap.Int("user_id", *request.UserID), zap.Any("request", request))

	existingUser, err := s.userQueryRepository.FindByIDWithPassword(ctx, *request.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,

			zap.Int("user_id", *request.UserID),
		)
	}

	if request.Email != "" && request.Email != existingUser.Email {
		duplicateUser, _ := s.userQueryRepository.FindByEmail(ctx, request.Email)
		if duplicateUser != nil {
			status = "error"
			return errorhandler.HandleError[*db.User](
				s.logger,
				user_errors.ErrUserEmailAlready,
				method,
				span,
				zap.String("email", request.Email),
			)
		}
		existingUser.Email = request.Email
	}

	if request.Password != "" {
		hash, err := s.hashing.HashPassword(request.Password)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.User](
				s.logger,
				user_errors.ErrUserPassword,
				method,
				span,
			)
		}
		existingUser.Password = hash
	}

	res, err := s.userCommandRepository.Update(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrFailedUpdateUser,
			method,
			span,

			zap.Int("user_id", *request.UserID),
		)
	}

	logSuccess("Successfully updated user", zap.Int("user_id", int(res.UserID)))

	return res, nil
}

func (s *userCommandService) Trash(ctx context.Context, user_id int) (*db.TrashUserRow, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", user_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Trashing user", zap.Int("user_id", user_id))

	res, err := s.userCommandRepository.Trash(ctx, user_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.TrashUserRow](
			s.logger,
			user_errors.ErrFailedTrashedUser,
			method,
			span,

			zap.Int("user_id", user_id),
		)
	}

	logSuccess("Successfully trashed user", zap.Int("user_id", user_id))

	return res, nil
}

func (s *userCommandService) Restore(ctx context.Context, user_id int) (*db.RestoreUserRow, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", user_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring user", zap.Int("user_id", user_id))

	res, err := s.userCommandRepository.Restore(ctx, user_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.RestoreUserRow](
			s.logger,
			user_errors.ErrFailedRestoreUser,
			method,
			span,

			zap.Int("user_id", user_id),
		)
	}

	logSuccess("Successfully restored user", zap.Int("user_id", user_id))

	return res, nil
}

func (s *userCommandService) DeletePermanent(ctx context.Context, user_id int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", user_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Deleting user permanently", zap.Int("user_id", user_id))

	_, err := s.userCommandRepository.DeletePermanent(ctx, user_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			user_errors.ErrFailedDeletePermanent,
			method,
			span,

			zap.Int("user_id", user_id),
		)
	}

	logSuccess("Successfully deleted user permanently", zap.Int("user_id", user_id))

	return true, nil
}

func (s *userCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring all users")

	_, err := s.userCommandRepository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			user_errors.ErrFailedRestoreAll,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all users")

	return true, nil
}

func (s *userCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Permanently deleting all users")

	_, err := s.userCommandRepository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			user_errors.ErrFailedDeleteAll,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all users permanently")

	return true, nil
}
func (s *userCommandService) UpdateIsVerified(ctx context.Context, user_id int, is_verified bool) (*db.User, error) {
	const method = "UpdateIsVerified"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", user_id),
		attribute.Bool("is_verified", is_verified))

	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.UpdateIsVerified(ctx, user_id, is_verified)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.User](
			s.logger,
			err,
			method,
			span,
			zap.Int("user_id", user_id),
		)
	}

	user := &db.User{
		UserID:     res.UserID,
		Firstname:  res.Firstname,
		Lastname:   res.Lastname,
		Email:      res.Email,
		IsVerified: &is_verified,
		CreatedAt:  res.CreatedAt,
		UpdatedAt:  res.UpdatedAt,
	}

	logSuccess("Successfully updated verification status", zap.Int("user_id", user_id))

	return user, nil
}

func (s *userCommandService) UpdatePassword(ctx context.Context, user_id int, password string) (*db.User, error) {
	const method = "UpdatePassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", user_id))

	defer func() {
		end(status)
	}()

	hash, err := s.hashing.HashPassword(password)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUserPassword,
			method,
			span,
		)
	}

	res, err := s.userCommandRepository.UpdatePassword(ctx, user_id, hash)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.User](
			s.logger,
			err,
			method,
			span,
			zap.Int("user_id", user_id),
		)
	}

	user := &db.User{
		UserID:    res.UserID,
		Firstname: res.Firstname,
		Lastname:  res.Lastname,
		Email:     res.Email,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}

	logSuccess("Successfully updated password", zap.Int("user_id", user_id))

	return user, nil
}
