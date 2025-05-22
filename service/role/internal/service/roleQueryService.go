package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type roleQueryService struct {
	ctx             context.Context
	trace           trace.Tracer
	roleQuery       repository.RoleQueryRepository
	logger          logger.LoggerInterface
	mapping         response_service.RoleResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewRoleQueryService(ctx context.Context, roleQuery repository.RoleQueryRepository, logger logger.LoggerInterface, mapping response_service.RoleResponseMapper) *roleQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "role_query_service_request_total",
			Help: "Total number of requests to the RoleQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "role_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the RoleQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &roleQueryService{
		ctx:             ctx,
		trace:           otel.Tracer("role-query-service"),
		roleQuery:       roleQuery,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *roleQueryService) FindAll(req *requests.FindAllRole) ([]*response.RoleResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, startTime)
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

	s.logger.Debug("Fetching role",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.roleQuery.FindAllRoles(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_ROLE")

		s.logger.Error("Failed to fetch role",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		status = "failed"
		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch role")
		status = "failed_to_fetch_role"

		return nil, nil, role_errors.ErrFailedFindAll
	}

	s.logger.Debug("Successfully fetched role",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	so := s.mapping.ToRolesResponse(res)

	return so, totalRecords, nil
}

func (s *roleQueryService) FindById(id int) (*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
	)

	s.logger.Debug("Fetching role by ID", zap.Int("id", id))

	res, err := s.roleQuery.FindById(id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ROLE_BY_ID")

		s.logger.Error("Failed to retrieve role by ID",
			zap.Int("id", id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve role by ID")
		status = "failed_to_retrieve_role_by_id"

		return nil, role_errors.ErrRoleNotFoundRes
	}

	s.logger.Debug("Successfully fetched role", zap.Int("id", id))

	so := s.mapping.ToRoleResponse(res)

	return so, nil
}

func (s *roleQueryService) FindByUserId(id int) ([]*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByUserId", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByUserId")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
	)

	s.logger.Debug("Fetching role by user ID", zap.Int("id", id))

	res, err := s.roleQuery.FindByUserId(id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ROLE_BY_USER_ID")

		s.logger.Error("Failed to retrieve role by user ID",
			zap.Int("id", id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve role by user ID")
		status = "failed_to_retrieve_role_by_user_id"

		return nil, role_errors.ErrRoleNotFoundRes
	}

	s.logger.Debug("Successfully fetched role by user ID", zap.Int("id", id))

	so := s.mapping.ToRolesResponse(res)

	return so, nil
}

func (s *roleQueryService) FindByActiveRole(req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActiveRole", status, startTime)
	}()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	_, span := s.trace.Start(s.ctx, "FindByActiveRole")
	defer span.End()

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching active role",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.roleQuery.FindByActiveRole(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_ROLE")

		s.logger.Error("Failed to fetch active role",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch active role")
		status = "failed_to_fetch_active_role"

		return nil, nil, role_errors.ErrFailedFindActive
	}

	s.logger.Debug("Successfully fetched active role",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	so := s.mapping.ToRolesResponseDeleteAt(res)

	return so, totalRecords, nil
}

func (s *roleQueryService) FindByTrashedRole(req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashedRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashedRole")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching trashed role",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.roleQuery.FindByTrashedRole(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_ROLE")

		s.logger.Error("Failed to fetch trashed role",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch trashed role")
		status = "failed_to_fetch_trashed_role"

		return nil, nil, role_errors.ErrFailedFindTrashed
	}

	s.logger.Debug("Successfully fetched trashed role",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	so := s.mapping.ToRolesResponseDeleteAt(res)

	return so, totalRecords, nil
}
func (s *roleQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
