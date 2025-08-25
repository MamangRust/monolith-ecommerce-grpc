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
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type roleCommandService struct {
	errorhandler    errorhandler.RoleCommandErrorHandler
	mencache        mencache.RoleCommandCache
	trace           trace.Tracer
	roleCommand     repository.RoleCommandRepository
	logger          logger.LoggerInterface
	mapping         response_service.RoleResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewRoleCommandService(errorhandler errorhandler.RoleCommandErrorHandler,
	mencache mencache.RoleCommandCache, roleCommand repository.RoleCommandRepository, logger logger.LoggerInterface, mapping response_service.RoleResponseMapper) *roleCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "role_command_service_request_total",
			Help: "Total number of requests to the RoleCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "role_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the RoleCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &roleCommandService{
		errorhandler:    errorhandler,
		mencache:        mencache,
		trace:           otel.Tracer("role-command-service"),
		roleCommand:     roleCommand,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *roleCommandService) CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "CreateRole"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.String("name", request.Name))

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.CreateRole(ctx, request)

	if err != nil {
		return s.errorhandler.HandleCreateRoleError(err, method, "FAILED_CREATE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	logSuccess("Successfully created role", zap.Int("role.id", role.ID), zap.Bool("success", true))

	return so, nil
}

func (s *roleCommandService) UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "UpdateRole"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("id", *request.ID))

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.UpdateRole(ctx, request)
	if err != nil {
		return s.errorhandler.HandleUpdateRoleError(err, method, "FAILED_UPDATE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	s.mencache.DeleteCachedRole(ctx, *request.ID)

	logSuccess("Successfully updated role", zap.Int("role.id", role.ID), zap.Bool("success", true))

	return so, nil
}

func (s *roleCommandService) TrashedRole(ctx context.Context, id int) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "TrashedRole"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.TrashedRole(ctx, id)

	if err != nil {
		return s.errorhandler.HandleTrashedRoleError(err, method, "FAILED_TRASH_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	s.mencache.DeleteCachedRole(ctx, id)

	logSuccess("Successfully trashed role", zap.Int("role.id", role.ID), zap.Bool("success", true))

	return so, nil
}

func (s *roleCommandService) RestoreRole(ctx context.Context, id int) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "RestoreRole"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.RestoreRole(ctx, id)

	if err != nil {
		return s.errorhandler.HandleRestoreRoleError(err, method, "FAILED_RESTORE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	s.mencache.DeleteCachedRole(ctx, id)

	logSuccess("Successfully restored role", zap.Int("role.id", role.ID), zap.Bool("success", true))

	return so, nil
}

func (s *roleCommandService) DeleteRolePermanent(ctx context.Context, id int) (bool, *response.ErrorResponse) {
	const method = "DeleteRolePermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	success, err := s.roleCommand.DeleteRolePermanent(ctx, id)
	if err != nil {
		return s.errorhandler.HandleDeleteRolePermanentError(err, method, "FAILED_DELETE_ROLE_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted role permanently", zap.Int("role.id", id), zap.Bool("success", success))

	return success, nil
}

func (s *roleCommandService) RestoreAllRole(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllRole"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.roleCommand.RestoreAllRole(ctx)
	if err != nil {
		return s.errorhandler.HandleRestoreAllRoleError(err, method, "FAILED_RESTORE_ALL_ROLE", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all roles", zap.Bool("success", success))

	return success, nil
}

func (s *roleCommandService) DeleteAllRolePermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllRolePermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.roleCommand.DeleteAllRolePermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllRolePermanentError(err, method, "FAILED_DELETE_ALL_ROLE_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all roles permanently", zap.Bool("success", success))

	return success, nil
}

func (s *roleCommandService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *roleCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
