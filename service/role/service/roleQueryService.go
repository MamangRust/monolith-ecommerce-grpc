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

type roleQueryService struct {
	observability  observability.TraceLoggerObservability
	cache          cache.RoleQueryCache
	roleRepository repository.RoleQueryRepository
	logger         logger.LoggerInterface
}

type RoleQueryServiceDeps struct {
	Observability  observability.TraceLoggerObservability
	Cache          cache.RoleQueryCache
	RoleRepository repository.RoleQueryRepository
	Logger         logger.LoggerInterface
}

func NewRoleQueryService(deps *RoleQueryServiceDeps) RoleQueryService {
	return &roleQueryService{
		observability:  deps.Observability,
		cache:          deps.Cache,
		roleRepository: deps.RoleRepository,
		logger:         deps.Logger,
	}
}

func (s *roleQueryService) FindAll(ctx context.Context, req *requests.FindAllRole) ([]*db.GetRolesRow, *int, error) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedRoles(ctx, req); found {
		logSuccess("Successfully retrieved all role records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, err := s.roleRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetRolesRow](
			s.logger,
			role_errors.ErrFindAllRoles,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedRoles(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched all roles",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return res, &totalCount, nil
}

func (s *roleQueryService) FindByID(ctx context.Context, id int) (*db.Role, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("role.id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedRoleById(ctx, id); found {
		logSuccess("Data found in cache", zap.Int("role.id", id))
		return data, nil
	}

	res, err := s.roleRepository.FindByID(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrRoleNotFound,
			method,
			span,
			zap.Int("role.id", id),
		)
	}

	s.cache.SetCachedRoleById(ctx, id, res)

	logSuccess("Successfully fetched role", zap.Int("role.id", id))

	return res, nil
}

func (s *roleQueryService) FindByName(ctx context.Context, name string) (*db.Role, error) {
	const method = "FindByName"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("role.name", name))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedRoleByName(ctx, name); found {
		logSuccess("Data found in cache", zap.String("role.name", name))
		return data, nil
	}

	res, err := s.roleRepository.FindByName(ctx, name)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrRoleNotFound,
			method,
			span,
			zap.String("role.name", name),
		)
	}

	s.cache.SetCachedRoleByName(ctx, name, res)

	logSuccess("Successfully fetched role", zap.String("role.name", name))

	return res, nil
}

func (s *roleQueryService) FindByUserId(ctx context.Context, id int) ([]*db.Role, error) {
	const method = "FindByUserId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedRoleByUserId(ctx, id); found {
		logSuccess("Data found in cache", zap.Int("user.id", id))
		return data, nil
	}

	res, err := s.roleRepository.FindByUserId(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.Role](
			s.logger,
			role_errors.ErrRoleNotFound,
			method,
			span,
			zap.Int("user.id", id),
		)
	}

	s.cache.SetCachedRoleByUserId(ctx, id, res)

	logSuccess("Successfully fetched role by user ID", zap.Int("user.id", id))

	return res, nil
}

func (s *roleQueryService) FindActive(ctx context.Context, req *requests.FindAllRole) ([]*db.GetActiveRolesRow, *int, error) {
	const method = "FindActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedRoleActive(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	res, err := s.roleRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetActiveRolesRow](
			s.logger,
			role_errors.ErrFindActiveRoles,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedRoleActive(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched active role", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return res, &totalCount, nil
}

func (s *roleQueryService) FindTrashed(ctx context.Context, req *requests.FindAllRole) ([]*db.GetTrashedRolesRow, *int, error) {
	const method = "FindTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedRoleTrashed(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	res, err := s.roleRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTrashedRolesRow](
			s.logger,
			role_errors.ErrFindTrashedRoles,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedRoleTrashed(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched trashed role", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return res, &totalCount, nil
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
