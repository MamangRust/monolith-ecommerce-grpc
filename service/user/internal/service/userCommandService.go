package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userCommandService struct {
	ctx                   context.Context
	trace                 trace.Tracer
	userQueryRepository   repository.UserQueryRepository
	userCommandRepository repository.UserCommandRepository
	roleRepository        repository.RoleRepository
	logger                logger.LoggerInterface
	mapping               response_service.UserResponseMapper
	hashing               hash.HashPassword
	requestCounter        *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
}

func NewUserCommandService(
	ctx context.Context,
	userQueryRepository repository.UserQueryRepository,
	userCommandRepository repository.UserCommandRepository,
	roleRepository repository.RoleRepository,
	logger logger.LoggerInterface,
	mapper response_service.UserResponseMapper,
	hashing hash.HashPassword,
) *userCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_command_service_requests_total",
			Help: "Total number of requests to the UserCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the UserCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &userCommandService{
		ctx:                   ctx,
		trace:                 otel.Tracer("user-command-service"),
		userQueryRepository:   userQueryRepository,
		userCommandRepository: userCommandRepository,
		roleRepository:        roleRepository,
		logger:                logger,
		mapping:               mapper,
		hashing:               hashing,
		requestCounter:        requestCounter,
		requestDuration:       requestDuration,
	}
}

func (s *userCommandService) CreateUser(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateUser", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("email", request.Email),
	)

	s.logger.Debug("Creating new user", zap.String("email", request.Email), zap.Any("request", request))

	existingUser, err := s.userQueryRepository.FindByEmail(request.Email)
	if err != nil {
		traceID := traceunic.GenerateTraceID("EMAIL_NOT_FOUND")

		s.logger.Error("Failed to find user by email", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		status = "user_not_found"

		return nil, user_errors.ErrUserEmailAlready

	} else if existingUser != nil {
		traceID := traceunic.GenerateTraceID("EMAIL_ALREADY_EXISTS")

		s.logger.Error("Email already exists", zap.String("trace_id", traceID))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Email already exists")
		status = "email_already_exists"

		return nil, user_errors.ErrUserEmailAlready
	}

	hash, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		traceID := traceunic.GenerateTraceID("PASSWORD_HASH_ERR")

		s.logger.Error("Failed to hash password", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to hash password")
		status = "password_hash_err"

		return nil, user_errors.ErrUserPassword
	}

	request.Password = hash

	const defaultRoleName = "Admin Access 1"
	role, err := s.roleRepository.FindByName(defaultRoleName)
	if err != nil || role == nil {
		traceID := traceunic.GenerateTraceID("ROLE_NOT_FOUND")

		s.logger.Error("Failed to find role by name", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Role not found")
		status = "role_not_found"

		return nil, role_errors.ErrRoleNotFoundRes
	}

	res, err := s.userCommandRepository.CreateUser(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_USER")

		s.logger.Error("Failed to create user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create user")
		status = "failed_create_user"

		return nil, user_errors.ErrFailedCreateUser
	}

	so := s.mapping.ToUserResponse(res)

	s.logger.Debug("Successfully created new user", zap.String("email", so.Email), zap.Int("user", so.ID))

	return so, nil
}

func (s *userCommandService) UpdateUser(request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateUser", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateUser")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", *request.UserID),
	)

	s.logger.Debug("Updating user", zap.Int("user_id", *request.UserID), zap.Any("request", request))

	existingUser, err := s.userQueryRepository.FindById(*request.UserID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("USER_NOT_FOUND")

		s.logger.Error("Failed to find user by id", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		status = "user_not_found"

		return nil, user_errors.ErrUserNotFoundRes
	}

	if request.Email != "" && request.Email != existingUser.Email {
		duplicateUser, _ := s.userQueryRepository.FindByEmail(request.Email)

		if duplicateUser != nil {
			traceID := traceunic.GenerateTraceID("EMAIL_ALREADY_EXISTS")

			s.logger.Error("Email already exists", zap.String("trace_id", traceID))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Email already exists")
			status = "email_already_exists"

			return nil, user_errors.ErrUserEmailAlready
		}

		existingUser.Email = request.Email
	}

	if request.Password != "" {
		hash, err := s.hashing.HashPassword(request.Password)
		if err != nil {
			traceID := traceunic.GenerateTraceID("PASSWORD_HASH_ERR")

			s.logger.Error("Failed to hash password", zap.String("trace_id", traceID), zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to hash password")
			status = "password_hash_err"

			return nil, user_errors.ErrUserPassword
		}
		existingUser.Password = hash
	}

	const defaultRoleName = "Admin Access 1"
	role, err := s.roleRepository.FindByName(defaultRoleName)
	if err != nil || role == nil {
		traceID := traceunic.GenerateTraceID("ROLE_NOT_FOUND")

		s.logger.Error("Failed to find role by name", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Role not found")
		status = "role_not_found"

		return nil, role_errors.ErrRoleNotFoundRes
	}

	res, err := s.userCommandRepository.UpdateUser(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_USER")

		s.logger.Error("Failed to update user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update user")
		status = "failed_update_user"

		return nil, user_errors.ErrFailedUpdateUser
	}

	so := s.mapping.ToUserResponse(res)

	s.logger.Debug("Successfully updated user", zap.Int("user_id", so.ID))

	return so, nil
}

func (s *userCommandService) TrashedUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedUser", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedUser")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", user_id),
	)

	s.logger.Debug("Trashing user", zap.Int("user_id", user_id))

	res, err := s.userCommandRepository.TrashedUser(user_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASHED_USER")

		s.logger.Error("Failed to trashed user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trashed user")
		status = "failed_trashed_user"

		return nil, user_errors.ErrFailedTrashedUser
	}

	so := s.mapping.ToUserResponseDeleteAt(res)

	s.logger.Debug("Successfully trashed user", zap.Int("user_id", user_id))

	return so, nil
}

func (s *userCommandService) RestoreUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreUser", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreUser")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", user_id),
	)

	s.logger.Debug("Restoring user", zap.Int("user_id", user_id))

	res, err := s.userCommandRepository.RestoreUser(user_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_USER")

		s.logger.Error("Failed to restore user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore user")
		status = "failed_restore_user"

		return nil, user_errors.ErrFailedRestoreUser
	}

	so := s.mapping.ToUserResponseDeleteAt(res)

	s.logger.Debug("Successfully restored user", zap.Int("user_id", user_id))

	return so, nil
}

func (s *userCommandService) DeleteUserPermanent(user_id int) (bool, *response.ErrorResponse) {
	start := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("DeleteUserPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteUserPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", user_id),
	)

	s.logger.Debug("Deleting user permanently", zap.Int("user_id", user_id))

	_, err := s.userCommandRepository.DeleteUserPermanent(user_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_USER_PERMANENT")

		s.logger.Error("Failed to permanently delete user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete user")
		status = "failed_delete_user_permanent"

		return false, user_errors.ErrFailedDeletePermanent
	}

	s.logger.Debug("Successfully deleted user permanently", zap.Int("user_id", user_id))

	return true, nil
}

func (s *userCommandService) RestoreAllUser() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllUser", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllUser")
	defer span.End()

	s.logger.Debug("Restoring all users")

	_, err := s.userCommandRepository.RestoreAllUser()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_USER")

		s.logger.Error("Failed to restore all users", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all users")
		status = "failed_restore_all_user"

		return false, user_errors.ErrFailedRestoreAll
	}

	s.logger.Debug("Successfully restored all users")

	return true, nil
}

func (s *userCommandService) DeleteAllUserPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllUserPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllUserPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all users")

	_, err := s.userCommandRepository.DeleteAllUserPermanent()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_USER_PERMANENT")

		s.logger.Error("Failed to permanently delete all users", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all users")
		status = "failed_delete_all_user_permanent"

		return false, user_errors.ErrFailedDeleteAll
	}

	s.logger.Debug("Successfully deleted all users permanently")

	return true, nil
}

func (s *userCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
