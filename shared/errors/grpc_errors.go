package errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGrpcError(err error) error {
	if err == nil {
		return nil
	}

	var apiErr *AppError
	if !errors.As(err, &apiErr) {
		return status.Error(codes.Internal, "Internal server error")
	}

	grpcCode := httpToGrpcCode(apiErr.Code)

	st := status.New(grpcCode, apiErr.Message)

	detail := &pb.ErrorResponse{
		Status:  apiErr.Type.String(),
		Message: apiErr.Message,
		Type:    string(apiErr.Type),
		Code:    int32(apiErr.Code),
	}

	stWithDetails, err := st.WithDetails(detail)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}

func ParseGrpcError(err error) *AppError {
	if err == nil {
		return nil
	}
	fmt.Printf("DEBUG: Parsing gRPC error: %v\n", err)

	st, ok := status.FromError(err)
	if !ok {
		return ErrInternal.WithInternal(err)
	}

	for _, detail := range st.Details() {
		if res, ok := detail.(*pb.ErrorResponse); ok {
			return &AppError{
				Type:    ErrorType(res.Type),
				Code:    int(res.Code),
				Message: res.Message,
			}
		}
	}

	// Fallback to gRPC status code
	return &AppError{
		Type:    grpcToErrorType(st.Code()),
		Code:    grpcToHttpCode(st.Code()),
		Message: st.Message(),
	}
}

func (t ErrorType) String() string {
	return string(t)
}

func httpToGrpcCode(code int) codes.Code {
	switch code {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	default:
		return codes.Internal
	}
}

func grpcToHttpCode(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

func grpcToErrorType(code codes.Code) ErrorType {
	switch code {
	case codes.NotFound:
		return ErrorTypeNotFound
	case codes.InvalidArgument:
		return ErrorTypeBadRequest
	case codes.AlreadyExists:
		return ErrorTypeConflict
	case codes.PermissionDenied:
		return ErrorTypeForbidden
	case codes.Unauthenticated:
		return ErrorTypeUnauthorized
	case codes.DeadlineExceeded:
		return ErrorTypeTimeout
	default:
		return ErrorTypeInternal
	}
}

func NewGrpcError(message string, httpCode int) error {
	grpcCode := httpToGrpcCode(httpCode)

	st := status.New(grpcCode, message)

	detail := &pb.ErrorResponse{
		Status:  http.StatusText(httpCode),
		Message: message,
		Code:    int32(httpCode),
	}

	stWithDetails, err := st.WithDetails(detail)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}

