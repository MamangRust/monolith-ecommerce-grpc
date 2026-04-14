package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/internal/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type roleCommandService struct {
	observability  observability.TraceLoggerObservability
	cache          cache.RoleCommandCache
	roleRepository repository.RoleCommandRepository
	logger         logger.LoggerInterface
}

type RoleCommandServiceDeps struct {
	Observability  observability.TraceLoggerObservability
	Cache          cache.RoleCommandCache
	RoleRepository repository.RoleCommandRepository
	Logger         logger.LoggerInterface
}

func NewRoleCommandService(deps *RoleCommandServiceDeps) RoleCommandService {
	return &roleCommandService{
		observability:  deps.Observability,
		cache:          deps.Cache,
		roleRepository: deps.RoleRepository,
		logger:         deps.Logger,
	}
}

func (s *roleCommandService) CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*db.Role, error) {
	const method = "CreateRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("name", request.Name))

	defer func() {
		end(status)
	}()

	role, err := s.roleRepository.CreateRole(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrCreateRole,
			method,
			span,
			zap.String("name", request.Name),
		)
	}

	logSuccess("Successfully created role", zap.Int32("role.id", role.RoleID))

	return role, nil
}

func (s *roleCommandService) UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*db.Role, error) {
	const method = "UpdateRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("id", *request.ID))

	defer func() {
		end(status)
	}()

	role, err := s.roleRepository.UpdateRole(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrUpdateRole,
			method,
			span,
			zap.Int("role.id", *request.ID),
		)
	}

	logSuccess("Successfully updated role", zap.Int32("role.id", role.RoleID))

	return role, nil
}

func (s *roleCommandService) TrashedRole(ctx context.Context, id int) (*db.Role, error) {
	const method = "TrashedRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	role, err := s.roleRepository.TrashedRole(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrTrashedRole,
			method,
			span,
			zap.Int("role.id", id),
		)
	}

	logSuccess("Successfully trashed role", zap.Int32("role.id", role.RoleID))

	return role, nil
}

func (s *roleCommandService) RestoreRole(ctx context.Context, id int) (*db.Role, error) {
	const method = "RestoreRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	role, err := s.roleRepository.RestoreRole(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrRestoreRole,
			method,
			span,
			zap.Int("role.id", id),
		)
	}

	logSuccess("Successfully restored role", zap.Int32("role.id", role.RoleID))

	return role, nil
}

func (s *roleCommandService) DeleteRolePermanent(ctx context.Context, id int) (bool, error) {
	const method = "DeleteRolePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	success, err := s.roleRepository.DeleteRolePermanent(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			role_errors.ErrDeleteRolePermanent,
			method,
			span,
			zap.Int("role.id", id),
		)
	}

	logSuccess("Successfully deleted role permanently", zap.Int("role.id", id))

	return success, nil
}

func (s *roleCommandService) RestoreAllRole(ctx context.Context) (bool, error) {
	const method = "RestoreAllRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.roleRepository.RestoreAllRole(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			role_errors.ErrRestoreAllRoles,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all roles")

	return success, nil
}

func (s *roleCommandService) DeleteAllRolePermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllRolePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.roleRepository.DeleteAllRolePermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			role_errors.ErrDeleteAllRoles,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all roles permanently")

	return success, nil
}
