# 📦 Package `userhandler`

**Source Path:** `apigateway/internal/handler/user`

## 🧩 Types

### `DepsUser`

```go
type DepsUser struct {
	Client *grpc.ClientConn
	E *echo.Echo
	Logger logger.LoggerInterface
}
```

### `userCommandHandleApi`

```go
type userCommandHandleApi struct {
	client pb.UserCommandServiceClient
	logger logger.LoggerInterface
	mapper apimapper.UserCommandResponseMapper
	trace trace.Tracer
	requestCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}
```

#### Methods

##### `Create`

@Security Bearer
Create handles the creation of a new user.
@Summary Create a new user
@Tags User Command
@Description Create a new user with the provided details
@Accept json
@Produce json
@Param request body requests.CreateUserRequest true "Create user request"
@Success 200 {object} response.ApiResponseUser "Successfully created user"
@Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
@Failure 500 {object} response.ErrorResponse "Failed to create user"
@Router /api/user-command/create [post]

```go
func (h *userCommandHandleApi) Create(c echo.Context) error
```

##### `DeleteAllUserPermanent`

@Security Bearer
DeleteUserPermanent permanently deletes a user record by its ID.
@Summary Permanently delete a user
@Tags User Command
@Description Permanently delete a user record by its ID.
@Accept json
@Produce json
@Param id path int true "User ID"
@Success 200 {object} response.ApiResponseUserDelete "Successfully deleted user record permanently"
@Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
@Failure 500 {object} response.ErrorResponse "Failed to delete user:"
@Router /api/user-command/delete/all [post]

```go
func (h *userCommandHandleApi) DeleteAllUserPermanent(c echo.Context) error
```

##### `DeleteUserPermanent`

@Security Bearer
DeleteUserPermanent permanently deletes a user record by its ID.
@Summary Permanently delete a user
@Tags User Command
@Description Permanently delete a user record by its ID.
@Accept json
@Produce json
@Param id path int true "User ID"
@Success 200 {object} response.ApiResponseUserDelete "Successfully deleted user record permanently"
@Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
@Failure 500 {object} response.ErrorResponse "Failed to delete user:"
@Router /api/user-command/delete/{id} [delete]

```go
func (h *userCommandHandleApi) DeleteUserPermanent(c echo.Context) error
```

##### `RestoreAllUser`

@Security Bearer
RestoreUser restores a user record from the trash by its ID.
@Summary Restore a trashed user
@Tags User Command
@Description Restore a trashed user record by its ID.
@Accept json
@Produce json
@Param id path int true "User ID"
@Success 200 {object} response.ApiResponseUserAll "Successfully restored user all"
@Failure 400 {object} response.ErrorResponse "Invalid user ID"
@Failure 500 {object} response.ErrorResponse "Failed to restore user"
@Router /api/user-command/restore/all [post]

```go
func (h *userCommandHandleApi) RestoreAllUser(c echo.Context) error
```

##### `RestoreUser`

@Security Bearer
RestoreUser restores a user record from the trash by its ID.
@Summary Restore a trashed user
@Tags User Command
@Description Restore a trashed user record by its ID.
@Accept json
@Produce json
@Param id path int true "User ID"
@Success 200 {object} response.ApiResponseUser "Successfully restored user"
@Failure 400 {object} response.ErrorResponse "Invalid user ID"
@Failure 500 {object} response.ErrorResponse "Failed to restore user"
@Router /api/user-command/restore/{id} [post]

```go
func (h *userCommandHandleApi) RestoreUser(c echo.Context) error
```

##### `TrashedUser`

@Security Bearer
TrashedUser retrieves a trashed user record by its ID.
@Summary Retrieve a trashed user
@Tags User Command
@Description Retrieve a trashed user record by its ID.
@Accept json
@Produce json
@Param id path int true "User ID"
@Success 200 {object} response.ApiResponseUser "Successfully retrieved trashed user"
@Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed user"
@Router /api/user-command/trashed/{id} [post]

```go
func (h *userCommandHandleApi) TrashedUser(c echo.Context) error
```

##### `Update`

@Security Bearer
Update handles the update of an existing user record.
@Summary Update an existing user
@Tags User Command
@Description Update an existing user record with the provided details
@Accept json
@Produce json
@Param id path int true "User ID"
@Param UpdateUserRequest body requests.UpdateUserRequest true "Update user request"
@Success 200 {object} response.ApiResponseUser "Successfully updated user"
@Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
@Failure 500 {object} response.ErrorResponse "Failed to update user"
@Router /api/user-command/update/{id} [post]

```go
func (h *userCommandHandleApi) Update(c echo.Context) error
```

##### `recordMetrics`

```go
func (s *userCommandHandleApi) recordMetrics(method string, status string, start time.Time)
```

##### `startTracingAndLogging`

```go
func (s *userCommandHandleApi) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (end func(), logSuccess func(string, ...zap.Field), logError func(string, error, ...zap.Field))
```

### `userCommandHandleParams`

```go
type userCommandHandleParams struct {
	client pb.UserCommandServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.UserCommandResponseMapper
}
```

### `userQueryHandleApi`

```go
type userQueryHandleApi struct {
	client pb.UserQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.UserQueryResponseMapper
	trace trace.Tracer
	requestCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}
```

#### Methods

##### `FindAllUser`

@Security Bearer
@Summary Find all users
@Tags User Command
@Description Retrieve a list of all users
@Accept json
@Produce json
@Param page query int false "Page number" default(1)
@Param page_size query int false "Number of items per page" default(10)
@Param search query string false "Search query"
@Success 200 {object} response.ApiResponsePaginationUser "List of users"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
@Router /api/user-query [get]

```go
func (h *userQueryHandleApi) FindAllUser(c echo.Context) error
```

##### `FindByActive`

@Security Bearer
@Summary Retrieve active users
@Tags User Command
@Description Retrieve a list of active users
@Accept json
@Produce json
@Param page query int false "Page number" default(1)
@Param page_size query int false "Number of items per page" default(10)
@Param search query string false "Search query"
@Success 200 {object} response.ApiResponsesUser "List of active users"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
@Router /api/user-query/active [get]

```go
func (h *userQueryHandleApi) FindByActive(c echo.Context) error
```

##### `FindById`

@Security Bearer
@Summary Find user by ID
@Tags User Command
@Description Retrieve a user by ID
@Accept json
@Produce json
@Param id path int true "User ID"
@Success 200 {object} response.ApiResponseUser "User data"
@Failure 400 {object} response.ErrorResponse "Invalid user ID"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
@Router /api/user-query/{id} [get]

```go
func (h *userQueryHandleApi) FindById(c echo.Context) error
```

##### `FindByTrashed`

@Security Bearer
FindByTrashed retrieves a list of trashed user records.
@Summary Retrieve trashed users
@Tags User Command
@Description Retrieve a list of trashed user records
@Accept json
@Produce json
@Param page query int false "Page number" default(1)
@Param page_size query int false "Number of items per page" default(10)
@Param search query string false "Search query"
@Success 200 {object} response.ApiResponsesUser "List of trashed user data"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
@Router /api/user-query/trashed [get]

```go
func (h *userQueryHandleApi) FindByTrashed(c echo.Context) error
```

##### `recordMetrics`

```go
func (s *userQueryHandleApi) recordMetrics(method string, status string, start time.Time)
```

##### `startTracingAndLogging`

```go
func (s *userQueryHandleApi) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (end func(), logSuccess func(string, ...zap.Field), logError func(string, error, ...zap.Field))
```

### `userQueryHandleParams`

```go
type userQueryHandleParams struct {
	client pb.UserQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.UserQueryResponseMapper
}
```

## 🚀 Functions

### `RegisterUserHandler`

```go
func RegisterUserHandler(deps *DepsUser)
```

### `setupUserCommandHandler`

```go
func setupUserCommandHandler(deps *DepsUser, mapper apimapper.UserCommandResponseMapper) func()
```

### `setupUserQueryHandler`

```go
func setupUserQueryHandler(deps *DepsUser, mapper apimapper.UserQueryResponseMapper) func()
```

