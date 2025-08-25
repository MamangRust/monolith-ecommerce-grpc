package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-role/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	errorhandler    errorhandler.RoleQueryErrorHandler
	mencache        mencache.RoleQueryCache
	trace           trace.Tracer
	roleQuery       repository.RoleQueryRepository
	logger          logger.LoggerInterface
	mapping         response_service.RoleResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewRoleQueryService(
	errorhandler errorhandler.RoleQueryErrorHandler,
	mencache mencache.RoleQueryCache,
	roleQuery repository.RoleQueryRepository, logger logger.LoggerInterface, mapping response_service.RoleResponseMapper) *roleQueryService {
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &roleQueryService{
		errorhandler:    errorhandler,
		mencache:        mencache,
		trace:           otel.Tracer("role-query-service"),
		roleQuery:       roleQuery,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *roleQueryService) FindAll(ctx context.Context, req *requests.FindAllRole) ([]*response.RoleResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoles(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindAllRoles(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(
			err, "FindAll", "FAILED_FIND_ALL_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRolesResponse(res)

	s.mencache.SetCachedRoles(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched roles", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *roleQueryService) FindById(ctx context.Context, id int) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("role.id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedRoleById(ctx, id); found {
		logSuccess("Data found in cache", zap.Int("role.id", id))

		return data, nil
	}

	res, err := s.roleQuery.FindById(ctx, id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(
			err, method, "FAILED_FIND_ROLE_BY_ID", span, &status, role_errors.ErrRoleNotFoundRes, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(res)

	s.mencache.SetCachedRoleById(ctx, so)

	logSuccess("Successfully fetched role", zap.Int("role.id", id), zap.Bool("success", true))

	return so, nil
}

func (s *roleQueryService) FindByUserId(ctx context.Context, id int) ([]*response.RoleResponse, *response.ErrorResponse) {
	const method = "FindByUserId"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedRoleByUserId(ctx, id); found {
		logSuccess("Data found in cache", zap.Int("user.id", id))

		return data, nil
	}

	res, err := s.roleQuery.FindByUserId(ctx, id)

	if err != nil {
		return s.errorhandler.HandleRepositoryListError(err, method, "FAILED_FIND_ROLE_BY_USER_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRolesResponse(res)

	s.mencache.SetCachedRoleByUserId(ctx, id, so)

	logSuccess("Successfully fetched role by user ID", zap.Int("user.id", id), zap.Bool("success", true))

	return so, nil
}

func (s *roleQueryService) FindByActiveRole(ctx context.Context, req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActiveRole"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoleActive(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindByActiveRole(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeletedError(err, method, "FAILED_FIND_BY_ACTIVE_ROLE", span, &status, role_errors.ErrRoleNotFoundRes, zap.Error(err))
	}

	so := s.mapping.ToRolesResponseDeleteAt(res)

	s.mencache.SetCachedRoleActive(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched active role", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *roleQueryService) FindByTrashedRole(ctx context.Context, req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashedRole"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoleTrashed(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindByTrashedRole(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeletedError(err, "FindByTrashedRole", "FAILED_FIND_BY_TRASHED_ROLE", span, &status, role_errors.ErrRoleNotFoundRes)
	}
	so := s.mapping.ToRolesResponseDeleteAt(res)

	s.mencache.SetCachedRoleTrashed(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed role", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *roleQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(ctx, method)

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

func (s *roleQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *roleQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
