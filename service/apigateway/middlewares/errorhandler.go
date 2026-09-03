package middlewares

import (
	"errors"

	sharederrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/labstack/echo/v4"
)

// RegisterErrorHandler installs a global error handler that maps shared AppError
// values to their proper HTTP status (400/401/403/404/409/429/503/504) instead of
// letting Echo's default handler collapse every non-HTTPError into 500.
//
// Handlers that return sharedErrors.ParseGrpcError(...) directly rely on this;
// routes wrapped in ApiHandler.Handle map errors themselves and are unaffected.
func RegisterErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var apiErr *sharederrors.AppError
		if errors.As(err, &apiErr) {
			response := sharederrors.ErrorResponse{
				Status:      "error",
				Message:     apiErr.Message,
				Type:        apiErr.Type,
				Code:        apiErr.Code,
				TraceID:     c.Response().Header().Get(echo.HeaderXRequestID),
				Retryable:   apiErr.Retryable,
				Validations: apiErr.Validations,
			}
			if writeErr := c.JSON(apiErr.Code, response); writeErr != nil {
				c.Logger().Error(writeErr)
			}
			return
		}
		// Preserve default handling for echo.HTTPError (auth middleware 401,
		// router 404, etc.) and any other error type.
		e.DefaultHTTPErrorHandler(err, c)
	}
}
