# 📦 Package `rolehandler`

**Source Path:** `apigateway/internal/handler/role`

## 🧩 Types

### `DepsRole`

```go
type DepsRole struct {
	Client *grpc.ClientConn
	Kafka *kafka.Kafka
	E *echo.Echo
	Logger logger.LoggerInterface
}
```

### `roleCommandHandleApi`

roleCommandHandleApi provides HTTP handlers for role-related operations.

This struct integrates the RoleService gRPC client, logging, Kafka event publishing,
response mapper, tracing, and Prometheus metrics for complete observability and functionality.

```go
type roleCommandHandleApi struct {
	kafka *kafka.Kafka
	role pb.RoleCommandServiceClient
	logger logger.LoggerInterface
	mapper apimapper.RoleCommandResponseMapper
	trace trace.Tracer
	requestCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}
```

#### Methods

##### `Create`

Create godoc.
@Summary Create a new role
@Tags Role
@Security Bearer
@Description Create a new role with the provided details.
@Accept json
@Produce json
@Param request body requests.CreateRoleRequest true "Role data"
@Success 200 {object} response.ApiResponseRole "Created role data"
@Failure 400 {object} response.ErrorResponse "Invalid request body"
@Failure 500 {object} response.ErrorResponse "Failed to create role"
@Router /api/role [post]

```go
func (h *roleCommandHandleApi) Create(c echo.Context) error
```

##### `DeleteAllPermanent`

DeleteAllPermanent godoc.
@Summary Permanently delete all roles
@Tags Role
@Security Bearer
@Description Permanently delete all roles.
@Accept json
@Produce json
@Success 200 {object} response.ApiResponseRoleAll "Permanently deleted roles data"
@Failure 500 {object} response.ErrorResponse "Failed to delete all roles permanently"
@Router /api/role/permanent/all [post]

```go
func (h *roleCommandHandleApi) DeleteAllPermanent(c echo.Context) error
```

##### `DeletePermanent`

DeletePermanent godoc.
@Summary Permanently delete a role
@Tags Role
@Security Bearer
@Description Permanently delete a role by its ID.
@Accept json
@Produce json
@Param id path int true "Role ID"
@Success 200 {object} response.ApiResponseRole "Permanently deleted role data"
@Failure 400 {object} response.ErrorResponse "Invalid role ID"
@Failure 500 {object} response.ErrorResponse "Failed to delete role permanently"
@Router /api/role/permanent/{id} [delete]

```go
func (h *roleCommandHandleApi) DeletePermanent(c echo.Context) error
```

##### `Restore`

Restore godoc.
@Summary Restore a soft-deleted role
@Tags Role
@Security Bearer
@Description Restore a soft-deleted role by its ID.
@Accept json
@Produce json
@Param id path int true "Role ID"
@Success 200 {object} response.ApiResponseRole "Restored role data"
@Failure 400 {object} response.ErrorResponse "Invalid role ID"
@Failure 500 {object} response.ErrorResponse "Failed to restore role"
@Router /api/role/restore/{id} [put]

```go
func (h *roleCommandHandleApi) Restore(c echo.Context) error
```

##### `RestoreAll`

RestoreAll godoc.
@Summary Restore all soft-deleted roles
@Tags Role
@Security Bearer
@Description Restore all soft-deleted roles.
@Accept json
@Produce json
@Success 200 {object} response.ApiResponseRoleAll "Restored roles data"
@Failure 500 {object} response.ErrorResponse "Failed to restore all roles"
@Router /api/role/restore/all [post]

```go
func (h *roleCommandHandleApi) RestoreAll(c echo.Context) error
```

##### `Trashed`

Trashed godoc.
@Summary Soft-delete a role
@Tags Role
@Security Bearer
@Description Soft-delete a role by its ID.
@Accept json
@Produce json
@Param id path int true "Role ID"
@Success 200 {object} response.ApiResponseRole "Soft-deleted role data"
@Failure 400 {object} response.ErrorResponse "Invalid role ID"
@Failure 500 {object} response.ErrorResponse "Failed to soft-delete role"
@Router /api/role/{id} [delete]

```go
func (h *roleCommandHandleApi) Trashed(c echo.Context) error
```

##### `Update`

Update godoc.
@Summary Update a role
@Tags Role
@Security Bearer
@Description Update an existing role with the provided details.
@Accept json
@Produce json
@Param id path int true "Role ID"
@Param request body requests.UpdateRoleRequest true "Role data"
@Success 200 {object} response.ApiResponseRole "Updated role data"
@Failure 400 {object} response.ErrorResponse "Invalid role ID or request body"
@Failure 500 {object} response.ErrorResponse "Failed to update role"
@Router /api/role/{id} [post]

```go
func (h *roleCommandHandleApi) Update(c echo.Context) error
```

##### `recordMetrics`

```go
func (s *roleCommandHandleApi) recordMetrics(method string, status string, start time.Time)
```

##### `startTracingAndLogging`

```go
func (s *roleCommandHandleApi) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (end func(), logSuccess func(string, ...zap.Field), logError func(string, error, ...zap.Field))
```

### `roleCommandHandleParams`

roleHandleParams contains the dependencies required to initialize and register
the role handler routes into the Echo router.

```go
type roleCommandHandleParams struct {
	client pb.RoleCommandServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.RoleCommandResponseMapper
	kafka *kafka.Kafka
}
```

### `roleQueryHandleParams`

```go
type roleQueryHandleParams struct {
	client pb.RoleQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.RoleQueryResponseMapper
	kafka *kafka.Kafka
}
```

### `roleQueryHandlerApi`

```go
type roleQueryHandlerApi struct {
	kafka *kafka.Kafka
	role pb.RoleQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.RoleQueryResponseMapper
	trace trace.Tracer
	requestCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}
```

#### Methods

##### `FindAll`

FindAll godoc.
@Summary Get all roles
@Tags Role
@Security Bearer
@Description Retrieve a paginated list of roles with optional search and pagination parameters.
@Accept json
@Produce json
@Param page query int false "Page number (default: 1)"
@Param page_size query int false "Number of items per page (default: 10)"
@Param search query string false "Search keyword"
@Success 200 {object} response.ApiResponsePaginationRole "List of roles"
@Failure 400 {object} response.ErrorResponse "Invalid query parameters"
@Failure 500 {object} response.ErrorResponse "Failed to fetch roles"
@Router /api/role [get]

```go
func (h *roleQueryHandlerApi) FindAll(c echo.Context) error
```

##### `FindByActive`

FindByActive godoc.
@Summary Get active roles
@Tags Role
@Security Bearer
@Description Retrieve a paginated list of active roles with optional search and pagination parameters.
@Accept json
@Produce json
@Param page query int false "Page number (default: 1)"
@Param page_size query int false "Number of items per page (default: 10)"
@Param search query string false "Search keyword"
@Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of active roles"
@Failure 400 {object} response.ErrorResponse "Invalid query parameters"
@Failure 500 {object} response.ErrorResponse "Failed to fetch active roles"
@Router /api/role/active [get]

```go
func (h *roleQueryHandlerApi) FindByActive(c echo.Context) error
```

##### `FindById`

FindById godoc.
@Summary Get a role by ID
@Tags Role
@Security Bearer
@Description Retrieve a role by its ID.
@Accept json
@Produce json
@Param id path int true "Role ID"
@Success 200 {object} response.ApiResponseRole "Role data"
@Failure 400 {object} response.ErrorResponse "Invalid role ID"
@Failure 500 {object} response.ErrorResponse "Failed to fetch role"
@Router /api/role/{id} [get]

```go
func (h *roleQueryHandlerApi) FindById(c echo.Context) error
```

##### `FindByTrashed`

FindByTrashed godoc.
@Summary Get trashed roles
@Tags Role
@Security Bearer
@Description Retrieve a paginated list of trashed roles with optional search and pagination parameters.
@Accept json
@Produce json
@Param page query int false "Page number (default: 1)"
@Param page_size query int false "Number of items per page (default: 10)"
@Param search query string false "Search keyword"
@Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of trashed roles"
@Failure 400 {object} response.ErrorResponse "Invalid query parameters"
@Failure 500 {object} response.ErrorResponse "Failed to fetch trashed roles"
@Router /api/role/trashed [get]

```go
func (h *roleQueryHandlerApi) FindByTrashed(c echo.Context) error
```

##### `FindByUserId`

FindByUserId godoc.
@Summary Get role by user ID
@Tags Role
@Security Bearer
@Description Retrieve a role by the associated user ID.
@Accept json
@Produce json
@Param user_id path int true "User ID"
@Success 200 {object} response.ApiResponseRole "Role data"
@Failure 400 {object} response.ErrorResponse "Invalid user ID"
@Failure 500 {object} response.ErrorResponse "Failed to fetch role by user ID"
@Router /api/role/user/{user_id} [get]

```go
func (h *roleQueryHandlerApi) FindByUserId(c echo.Context) error
```

##### `recordMetrics`

```go
func (s *roleQueryHandlerApi) recordMetrics(method string, status string, start time.Time)
```

##### `startTracingAndLogging`

```go
func (s *roleQueryHandlerApi) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (end func(), logSuccess func(string, ...zap.Field), logError func(string, error, ...zap.Field))
```

## 🚀 Functions

### `RegisterRoleHandler`

```go
func RegisterRoleHandler(deps *DepsRole)
```

### `setupRoleCommandHandler`

```go
func setupRoleCommandHandler(deps *DepsRole, mapper apimapper.RoleCommandResponseMapper) func()
```

### `setupRoleQueryHandler`

```go
func setupRoleQueryHandler(deps *DepsRole, mapper apimapper.RoleQueryResponseMapper) func()
```

