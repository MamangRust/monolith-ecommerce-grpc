package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-user/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userQueryService struct {
	ctx                 context.Context
	errorhandler        errorhandler.UserQueryError
	mencache            mencache.UserQueryCache
	trace               trace.Tracer
	userQueryRepository repository.UserQueryRepository
	logger              logger.LoggerInterface
	mapping             response_service.UserResponseMapper
	requestCounter      *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
}

func NewUserQueryService(ctx context.Context,
	errorhandler errorhandler.UserQueryError,
	mencache mencache.UserQueryCache,
	userQueryRepository repository.UserQueryRepository, logger logger.LoggerInterface, mapping response_service.UserResponseMapper) *userQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_query_service_request_total",
			Help: "Total number of requests to the UserQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the UserQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &userQueryService{
		ctx:                 ctx,
		errorhandler:        errorhandler,
		mencache:            mencache,
		trace:               otel.Tracer("user-query-service"),
		userQueryRepository: userQueryRepository,
		logger:              logger,
		mapping:             mapping,
		requestCounter:      requestCounter,
		requestDuration:     requestDuration,
	}
}

func (s *userQueryService) FindAll(req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUsersCache(req); found {
		logSuccess("Successfully fetched user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindAllUsers(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_USERS", span, &status, zap.Error(err))
	}

	userResponses := s.mapping.ToUsersResponse(users)

	s.mencache.SetCachedUsersCache(req, userResponses, totalRecords)

	logSuccess("Successfully fetched all user records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return userResponses, totalRecords, nil
}

func (s *userQueryService) FindByID(id int) (*response.UserResponse, *response.ErrorResponse) {
	const method = "FindByID"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	user, err := s.userQueryRepository.FindById(id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes, zap.Int("user.id", id), zap.Error(err))
	}

	userRes := s.mapping.ToUserResponse(user)

	s.mencache.SetCachedUserCache(userRes)

	logSuccess("Successfully fetched user record", zap.Int("user.id", id))

	return userRes, nil
}

func (s *userQueryService) FindByActive(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserActiveCache(req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByActive", "FAILED_FIND_ACTIVE_USERS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.mencache.SetCachedUserActiveCache(req, so, totalRecords)

	logSuccess("Successfully fetched active users", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *userQueryService) FindByTrashed(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserTrashedCache(req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByTrashed", "FAILED_FIND_TRASHED_USERS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.mencache.SetCachedUserTrashedCache(req, so, totalRecords)

	logSuccess("Successfully fetched trashed users", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *userQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

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

	return span, end, status, logSuccess
}

func (s *userQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *userQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
