package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/repository"
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
	observability      observability.TraceLoggerObservability
	cache              cache.RoleCommandCache
	roleRepository     repository.RoleCommandRepository
	userRoleRepository repository.UserRoleRepository
	logger             logger.LoggerInterface
}

type RoleCommandServiceDeps struct {
	Observability      observability.TraceLoggerObservability
	Cache              cache.RoleCommandCache
	RoleRepository     repository.RoleCommandRepository
	UserRoleRepository repository.UserRoleRepository
	Logger             logger.LoggerInterface
}

func NewRoleCommandService(deps *RoleCommandServiceDeps) RoleCommandService {
	return &roleCommandService{
		observability:      deps.Observability,
		cache:              deps.Cache,
		roleRepository:     deps.RoleRepository,
		userRoleRepository: deps.UserRoleRepository,
		logger:             deps.Logger,
	}
}

func (s *roleCommandService) Create(ctx context.Context, request *requests.CreateRoleRequest) (*db.Role, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("name", request.Name))

	defer func() {
		end(status)
	}()

	role, err := s.roleRepository.Create(ctx, request)
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

func (s *roleCommandService) Update(ctx context.Context, request *requests.UpdateRoleRequest) (*db.Role, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("id", *request.ID))

	defer func() {
		end(status)
	}()

	role, err := s.roleRepository.Update(ctx, request)
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

func (s *roleCommandService) Trash(ctx context.Context, id int) (*db.Role, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	role, err := s.roleRepository.Trash(ctx, id)
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

func (s *roleCommandService) Restore(ctx context.Context, id int) (*db.Role, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	role, err := s.roleRepository.Restore(ctx, id)
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

func (s *roleCommandService) DeletePermanent(ctx context.Context, id int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	success, err := s.roleRepository.DeletePermanent(ctx, id)
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

func (s *roleCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.roleRepository.RestoreAll(ctx)
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

func (s *roleCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.roleRepository.DeleteAll(ctx)
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

func (s *roleCommandService) AssignRoleToUser(ctx context.Context, request *requests.CreateUserRoleRequest) (*db.UserRole, error) {
	const method = "AssignRoleToUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", request.UserId),
		attribute.Int("role_id", request.RoleId))

	defer func() {
		end(status)
	}()

	userRole, err := s.userRoleRepository.AssignRoleToUser(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UserRole](
			s.logger,
			err,
			method,
			span,
			zap.Int("user_id", request.UserId),
			zap.Int("role_id", request.RoleId),
		)
	}

	logSuccess("Successfully assigned role to user", zap.Int("user_id", request.UserId), zap.Int("role_id", request.RoleId))

	return userRole, nil
}

func (s *roleCommandService) RemoveRoleFromUser(ctx context.Context, request *requests.RemoveUserRoleRequest) error {
	const method = "RemoveRoleFromUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", request.UserId),
		attribute.Int("role_id", request.RoleId))

	defer func() {
		end(status)
	}()

	err := s.userRoleRepository.RemoveRoleFromUser(ctx, request)
	if err != nil {
		status = "error"
		_, err := errorhandler.HandleError[any](
			s.logger,
			err,
			method,
			span,
			zap.Int("user_id", request.UserId),
			zap.Int("role_id", request.RoleId),
		)
		return err
	}

	logSuccess("Successfully removed role from user", zap.Int("user_id", request.UserId), zap.Int("role_id", request.RoleId))

	return nil
}
