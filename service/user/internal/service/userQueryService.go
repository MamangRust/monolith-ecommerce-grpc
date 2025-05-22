package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
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
	trace               trace.Tracer
	userQueryRepository repository.UserQueryRepository
	logger              logger.LoggerInterface
	mapping             response_service.UserResponseMapper
	requestCounter      *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
}

func NewUserQueryService(ctx context.Context, userQueryRepository repository.UserQueryRepository, logger logger.LoggerInterface, mapping response_service.UserResponseMapper) *userQueryService {
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
		trace:               otel.Tracer("user-query-service"),
		userQueryRepository: userQueryRepository,
		logger:              logger,
		mapping:             mapping,
		requestCounter:      requestCounter,
		requestDuration:     requestDuration,
	}
}

func (s *userQueryService) FindAll(req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAll")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching users",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.userQueryRepository.FindAllUsers(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL")

		s.logger.Error("Failed to fetch users",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch users")
		status = "failed_find_all"

		return nil, nil, user_errors.ErrFailedFindAll
	}

	userResponses := s.mapping.ToUsersResponse(users)

	s.logger.Debug("Successfully fetched user",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return userResponses, totalRecords, nil
}

func (s *userQueryService) FindByID(id int) (*response.UserResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByID", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByID")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", id),
	)

	s.logger.Debug("Fetching user by id", zap.Int("user_id", id))

	user, err := s.userQueryRepository.FindById(id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("USER_NOT_FOUND")

		s.logger.Error("Failed to find user by id", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		status = "user_not_found"

		return nil, user_errors.ErrUserNotFoundRes
	}

	so := s.mapping.ToUserResponse(user)

	s.logger.Debug("Successfully fetched user", zap.Int("user_id", id))

	return so, nil
}

func (s *userQueryService) FindByActive(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByActive")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching active user",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.userQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE")

		s.logger.Error("Failed to find active users", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find active users")
		status = "failed_find_active"

		return nil, nil, user_errors.ErrFailedFindActive
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.logger.Debug("Successfully fetched active user",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *userQueryService) FindByTrashed(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashed")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching trashed user",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.userQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED")

		s.logger.Error("Failed to find trashed users", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find trashed users")
		status = "failed_find_trashed"

		return nil, nil, user_errors.ErrFailedFindTrashed
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.logger.Debug("Successfully fetched trashed user",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *userQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
