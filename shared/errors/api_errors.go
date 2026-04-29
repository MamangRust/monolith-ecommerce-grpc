package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
)

func InvalidAccessToken() error {
	return fmt.Errorf("invalid access token")
}

type ApiHandler interface {
	Handle(method string, handler func(echo.Context) error) echo.HandlerFunc
	HandleApiErrorWithTracing(c echo.Context, err error, span trace.Span, method string) error
}

type apiHandler struct {
	observability observability.TraceLoggerObservability
	logger        logger.LoggerInterface
}

func NewApiHandler(observability observability.TraceLoggerObservability, logger logger.LoggerInterface) ApiHandler {
	return &apiHandler{
		observability: observability,
		logger:        logger,
	}
}

func (h *apiHandler) Handle(method string, handler func(echo.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
			c.Request().Context(),
			method,
			attribute.String("path", c.Request().URL.Path),
			attribute.String("method", c.Request().Method),
		)

		c.SetRequest(c.Request().WithContext(ctx))

		defer func() {
			end(status)
		}()

		err := handler(c)
		if err != nil {
			status = "error"
			return h.HandleApiErrorWithTracing(c, err, span, method)
		}

		logSuccess("Request completed successfully")
		return nil
	}
}

func (h *apiHandler) HandleApiErrorWithTracing(c echo.Context, err error, span trace.Span, method string) error {
	traceID := span.SpanContext().TraceID().String()

	h.logger.Error(
		fmt.Sprintf("API error in %s", method),
		zap.Error(err),
		zap.String("trace.id", traceID),
		zap.String("path", c.Request().URL.Path),
		zap.String("method", c.Request().Method),
	)

	span.SetAttributes(
		attribute.String("trace.id", traceID),
		attribute.String("error", err.Error()),
	)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	return HandleApiError(c, err, traceID)
}

type ErrorResponse struct {
	Status      string            `json:"status"`
	Message     string            `json:"message"`
	Type        ErrorType         `json:"type,omitempty"`
	Code        int               `json:"code"`
	TraceID     string            `json:"trace_id,omitempty"`
	Retryable   bool              `json:"retryable,omitempty"`
	Validations []ValidationError `json:"validations,omitempty"`
}

func HandleApiError(c echo.Context, err error, traceID string) error {
	if err == nil {
		return nil
	}

	var apiErr *AppError
	if errors.As(err, &apiErr) {
		response := ErrorResponse{
			Status:      "error",
			Message:     apiErr.Message,
			Type:        apiErr.Type,
			Code:        apiErr.Code,
			TraceID:     traceID,
			Retryable:   apiErr.Retryable,
			Validations: apiErr.Validations,
		}
		return c.JSON(apiErr.Code, response)
	}

	response := ErrorResponse{
		Status:  "error",
		Message: "An internal server error occurred",
		Type:    ErrorTypeInternal,
		Code:    http.StatusInternalServerError,
		TraceID: traceID,
	}
	return c.JSON(http.StatusInternalServerError, response)
}

func NewBadRequestError(message string) *AppError {
	return ErrBadRequest.WithMessage(message)
}

func NewNotFoundError(resource string) *AppError {
	return ErrNotFound.WithMessage(fmt.Sprintf("%s not found", resource))
}

func NewConflictError(message string) *AppError {
	return ErrConflict.WithMessage(message)
}

func NewInternalError(err error) *AppError {
	return ErrInternal.WithInternal(err)
}

func NewServiceUnavailableError(service string) *AppError {
	return ErrServiceUnavailable.WithMessage(fmt.Sprintf("%s is temporarily unavailable", service))
}

