package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-user/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	errorhandler          errorhandler.UserCommandError
	mencache              mencache.UserCommandCache
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
	errorhandler errorhandler.UserCommandError,
	mencache mencache.UserCommandCache,
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
		mencache:              mencache,
		errorhandler:          errorhandler,
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

func (s *userCommandService) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	const method = "CreateUser"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	existingUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER_BY_EMAIL", span, &status, user_errors.ErrUserEmailAlready, zap.String("email", request.Email), zap.Error(err))

	} else if existingUser != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER_BY_EMAIL", span, &status, user_errors.ErrUserEmailAlready, zap.String("email", request.Email), zap.Int("user.id", existingUser.ID), zap.Error(err))
	}

	hash, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		return errorhandler.HandleErrorPasswordOperation[*response.UserResponse](s.logger, err, method, "FAILED_HASH_PASSWORD", span, &status, user_errors.ErrUserPassword, zap.String("email", request.Email), zap.Error(err))
	}

	request.Password = hash

	const defaultRoleName = "Admin Access 1"

	role, err := s.roleRepository.FindByName(ctx, defaultRoleName)

	if err != nil || role == nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_ROLE", span, &status, role_errors.ErrRoleNotFoundRes, zap.String("name", defaultRoleName), zap.Error(err))
	}

	res, err := s.userCommandRepository.CreateUser(ctx, request)

	if err != nil {
		return s.errorhandler.HandleCreateUserError(err, method, "FAILED_CREATE_USER", span, &status, zap.String("email", request.Email), zap.Error(err))
	}

	so := s.mapping.ToUserResponse(res)

	logSuccess("Successfully created user", zap.Int("user.id", res.ID), zap.Bool("success", true))

	return so, nil
}

func (s *userCommandService) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	const method = "UpdateUser"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	existingUser, err := s.userQueryRepository.FindById(ctx, *request.UserID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes)
	}

	if request.Email != "" && request.Email != existingUser.Email {
		duplicateUser, _ := s.userQueryRepository.FindByEmail(ctx, request.Email)

		if duplicateUser != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_EMAIL_ALREADY", span, &status, user_errors.ErrUserEmailAlready, zap.String("email", request.Email), zap.Int("user.id", duplicateUser.ID), zap.Error(err))
		}

		existingUser.Email = request.Email
	}

	if request.Password != "" {
		hash, err := s.hashing.HashPassword(request.Password)
		if err != nil {
			return errorhandler.HandleErrorPasswordOperation[*response.UserResponse](s.logger, err, method, "FAILED_HASH_PASSWORD", span, &status, user_errors.ErrUserPassword, zap.Int("user.id", *request.UserID), zap.Error(err))
		}
		existingUser.Password = hash
	}

	const defaultRoleName = "Admin Access 1"

	role, err := s.roleRepository.FindByName(ctx, defaultRoleName)

	if err != nil || role == nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_ROLE", span, &status, role_errors.ErrRoleNotFoundRes)
	}

	res, err := s.userCommandRepository.UpdateUser(ctx, request)

	if err != nil {
		return s.errorhandler.HandleUpdateUserError(err, method, "FAILED_UPDATE_USER", span, &status, zap.Int("user.id", *request.UserID), zap.Error(err))
	}

	so := s.mapping.ToUserResponse(res)
	s.mencache.DeleteUserCache(ctx, so.ID)

	logSuccess("Successfully updated user", zap.Int("user.id", res.ID), zap.Bool("success", true))

	return so, nil
}

func (s *userCommandService) TrashedUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedUser"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.TrashedUser(ctx, user_id)

	if err != nil {
		return s.errorhandler.HandleTrashedUserError(err, method, "FAILED_TO_TRASH_USER", span, &status, zap.Int("user.id", user_id), zap.Error(err))
	}

	so := s.mapping.ToUserResponseDeleteAt(res)

	s.mencache.DeleteUserCache(ctx, so.ID)

	logSuccess("Successfully trashed user", zap.Int("user.id", res.ID), zap.Bool("success", true))

	return so, nil
}

func (s *userCommandService) RestoreUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreUser"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.RestoreUser(ctx, user_id)

	if err != nil {
		return s.errorhandler.HandleRestoreUserError(err, method, "FAILED_TO_RESTORE_USER", span, &status, zap.Int("user.id", user_id), zap.Error(err))
	}

	so := s.mapping.ToUserResponseDeleteAt(res)

	logSuccess("Successfully restored user", zap.Int("user.id", res.ID), zap.Bool("success", true))

	return so, nil
}

func (s *userCommandService) DeleteUserPermanent(ctx context.Context, user_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteUserPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.userCommandRepository.DeleteUserPermanent(ctx, user_id)

	if err != nil {
		return s.errorhandler.HandleDeleteUserError(err, method, "FAILED_TO_DELETE_USER", span, &status, zap.Int("user.id", user_id), zap.Error(err))
	}

	logSuccess("Successfully deleted user", zap.Int("user.id", user_id), zap.Bool("success", success))

	return true, nil
}

func (s *userCommandService) RestoreAllUser(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllUser"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.userCommandRepository.RestoreAllUser(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllUserError(err, method, "FAILED_RESTORE_ALL_USER", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all users", zap.Bool("success", success))

	return true, nil
}

func (s *userCommandService) DeleteAllUserPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllUserPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.userCommandRepository.DeleteAllUserPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllUserError(err, method, "FAILED_DELETE_ALL_USER", span, &status)
	}

	logSuccess("Successfully deleted all users", zap.Bool("success", success))

	return true, nil
}

func (s *userCommandService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
}

func (s *userCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
