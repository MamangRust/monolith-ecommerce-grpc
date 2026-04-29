package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type userQueryService struct {
	observability  observability.TraceLoggerObservability
	cache          cache.UserQueryCache
	userRepository repository.UserQueryRepository
	logger         logger.LoggerInterface
}

type UserQueryServiceDeps struct {
	Observability  observability.TraceLoggerObservability
	Cache          cache.UserQueryCache
	UserRepository repository.UserQueryRepository
	Logger         logger.LoggerInterface
}

func NewUserQueryService(deps *UserQueryServiceDeps) UserQueryService {
	return &userQueryService{
		observability:  deps.Observability,
		cache:          deps.Cache,
		userRepository: deps.UserRepository,
		logger:         deps.Logger,
	}
}

func (s *userQueryService) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, *int, error) {
	const method = "FindAll"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedUsersCache(ctx, req); found {
		logSuccess("Successfully retrieved all user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	users, err := s.userRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetUsersRow](
			s.logger,
			user_errors.ErrFailedFindAll,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(users) > 0 {
		totalCount = int(users[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedUsersCache(ctx, req, users, &totalCount)

	logSuccess("Successfully fetched user",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return users, &totalCount, nil
}

func (s *userQueryService) FindByID(ctx context.Context, id int) (*db.GetUserByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedUserCache(ctx, id); found {
		logSuccess("Successfully retrieved user record from cache", zap.Int("user.id", id))
		return data, nil
	}

	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetUserByIDRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,

			zap.Int("user_id", id),
		)
	}

	s.cache.SetCachedUserCache(ctx, user)

	logSuccess("Successfully fetched user", zap.Int("user_id", id))

	return user, nil
}

func (s *userQueryService) FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, *int, error) {
	const method = "FindActive"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedUserActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	users, err := s.userRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetUsersActiveRow](
			s.logger,
			user_errors.ErrFailedFindActive,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(users) > 0 {
		totalCount = int(users[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedUserActiveCache(ctx, req, users, &totalCount)

	logSuccess("Successfully fetched active user",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return users, &totalCount, nil
}

func (s *userQueryService) FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, *int, error) {
	const method = "FindTrashed"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedUserTrashedCache(ctx, req); found {
		return data, total, nil
	}

	users, err := s.userRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetUserTrashedRow](
			s.logger,
			user_errors.ErrFailedFindTrashed,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(users) > 0 {
		totalCount = int(users[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedUserTrashedCache(ctx, req, users, &totalCount)

	logSuccess("Successfully fetched trashed user",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return users, &totalCount, nil
}
func (s *userQueryService) FindByEmailWithPassword(ctx context.Context, email string) (*db.GetUserByEmailWithPasswordRow, error) {
	const method = "FindByEmailWithPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("email", email))

	defer func() {
		end(status)
	}()

	user, err := s.userRepository.FindByEmailWithPassword(ctx, email)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetUserByEmailWithPasswordRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,
			zap.String("email", email),
		)
	}

	logSuccess("Successfully fetched user by email", zap.String("email", email))

	return user, nil
}

func (s *userQueryService) FindByVerificationCode(ctx context.Context, code string) (*db.GetUserByVerificationCodeRow, error) {
	const method = "FindByVerificationCode"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("verification_code", code))

	defer func() {
		end(status)
	}()

	user, err := s.userRepository.FindByVerificationCode(ctx, code)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetUserByVerificationCodeRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,
			zap.String("verification_code", code),
		)
	}

	logSuccess("Successfully fetched user by verification code", zap.String("verification_code", code))

	return user, nil
}
