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

type roleCommandService struct {
	ctx             context.Context
	trace           trace.Tracer
	roleCommand     repository.RoleCommandRepository
	logger          logger.LoggerInterface
	mapping         response_service.RoleResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewRoleCommandService(ctx context.Context, roleCommand repository.RoleCommandRepository, logger logger.LoggerInterface, mapping response_service.RoleResponseMapper) *roleCommandService {
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &roleCommandService{
		ctx:             ctx,
		trace:           otel.Tracer("role-command-service"),
		roleCommand:     roleCommand,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *roleCommandService) CreateRole(request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateRole")
	defer span.End()

	span.SetAttributes(
		attribute.String("name", request.Name),
	)

	s.logger.Debug("Starting CreateRole process",
		zap.String("roleName", request.Name),
	)

	role, err := s.roleCommand.CreateRole(request)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_ROLE")

		s.logger.Error("Failed to create role record",
			zap.String("trace_id", traceID),
			zap.String("name", request.Name),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create role record")
		status = "failed_to_create_role_record"
		return nil, role_errors.ErrFailedCreateRole
	}

	so := s.mapping.ToRoleResponse(role)

	s.logger.Debug("CreateRole process completed",
		zap.String("roleName", request.Name),
		zap.Int("roleID", role.ID),
	)

	return so, nil
}

func (s *roleCommandService) UpdateRole(request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateRole")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", *request.ID),
		attribute.String("name", request.Name),
	)

	s.logger.Debug("Starting UpdateRole process",
		zap.Int("roleID", *request.ID),
		zap.String("newRoleName", request.Name),
	)

	role, err := s.roleCommand.UpdateRole(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_ROLE")

		s.logger.Error("Failed to update role record",
			zap.String("trace_id", traceID),
			zap.Int("role_id", *request.ID),
			zap.String("name", request.Name),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update role record")
		status = "failed_to_update_role_record"
		return nil, role_errors.ErrFailedUpdateRole
	}

	so := s.mapping.ToRoleResponse(role)

	s.logger.Debug("UpdateRole process completed",
		zap.Int("roleID", *request.ID),
		zap.String("newRoleName", request.Name),
	)

	return so, nil
}

func (s *roleCommandService) TrashedRole(id int) (*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedRole")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
	)

	s.logger.Debug("Starting TrashedRole process",
		zap.Int("roleID", id),
	)

	role, err := s.roleCommand.TrashedRole(id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASHED_ROLE")

		s.logger.Error("Failed to trashed role",
			zap.String("trace_id", traceID),
			zap.Int("role_id", id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trashed role")
		status = "failed_to_trashed_role"
		return nil, role_errors.ErrFailedTrashedRole
	}

	so := s.mapping.ToRoleResponse(role)

	s.logger.Debug("TrashedRole process completed",
		zap.Int("roleID", id),
	)

	return so, nil
}

func (s *roleCommandService) RestoreRole(id int) (*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreRole")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
	)

	s.logger.Debug("Starting RestoreRole process",
		zap.Int("roleID", id),
	)

	role, err := s.roleCommand.RestoreRole(id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ROLE")

		s.logger.Error("Failed to restore role",
			zap.String("trace_id", traceID),
			zap.Int("role_id", id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore role")
		status = "failed_to_restore_role"

		return nil, role_errors.ErrFailedRestoreRole
	}

	so := s.mapping.ToRoleResponse(role)

	s.logger.Debug("RestoreRole process completed",
		zap.Int("roleID", id),
	)

	return so, nil
}

func (s *roleCommandService) DeleteRolePermanent(id int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteRolePermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteRolePermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
	)

	s.logger.Debug("Starting DeleteRolePermanent process",
		zap.Int("roleID", id),
	)

	_, err := s.roleCommand.DeleteRolePermanent(id)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ROLE_PERMANENT")

		s.logger.Error("Failed to permanently delete role",
			zap.String("trace_id", traceID),
			zap.Int("role_id", id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete role")
		status = "failed_to_permanently_delete_role"

		return false, role_errors.ErrFailedDeletePermanent
	}

	s.logger.Debug("DeleteRolePermanent process completed",
		zap.Int("roleID", id),
	)

	return true, nil
}

func (s *roleCommandService) RestoreAllRole() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllRole")
	defer span.End()

	_, err := s.roleCommand.RestoreAllRole()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_ROLE")

		s.logger.Error("Failed to restore all trashed roles",
			zap.String("trace_id", traceID),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all trashed roles")

		status = "failed_to_restore_all_trashed_roles"
		return false, role_errors.ErrFailedRestoreAll
	}

	s.logger.Debug("Successfully restored all roles")
	return true, nil
}

func (s *roleCommandService) DeleteAllRolePermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllRolePermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllRolePermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all roles")

	_, err := s.roleCommand.DeleteAllRolePermanent()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_ROLE_PERMANENT")

		s.logger.Error("Failed to permanently delete all roles",
			zap.String("trace_id", traceID),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all roles")
		status = "failed_to_permanently_delete_all_roles"

		return false, role_errors.ErrFailedDeletePermanent
	}

	s.logger.Debug("Successfully deleted all roles permanently")
	return true, nil
}

func (s *roleCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
