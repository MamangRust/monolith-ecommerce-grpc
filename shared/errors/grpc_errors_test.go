package errors

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// roundTripCases asserts that an HTTP code survives a gRPC round trip:
// AppError(http) -> ToGrpcError -> ParseGrpcError -> same HTTP code.
var roundTripCases = []struct {
	name string
	code int
}{
	{"bad_request", http.StatusBadRequest},
	{"unauthorized", http.StatusUnauthorized},
	{"forbidden", http.StatusForbidden},
	{"not_found", http.StatusNotFound},
	{"conflict", http.StatusConflict},
	{"too_many_requests", http.StatusTooManyRequests},
	{"service_unavailable", http.StatusServiceUnavailable},
	{"gateway_timeout", http.StatusGatewayTimeout},
}

func TestRoundTripHTTPToGRPCAndBack(t *testing.T) {
	for _, tc := range roundTripCases {
		t.Run(tc.name, func(t *testing.T) {
			appErr := &AppError{
				Type:    ErrorType("TEST"),
				Code:    tc.code,
				Message: "test message",
			}

			grpcErr := ToGrpcError(appErr)
			if grpcErr == nil {
				t.Fatal("ToGrpcError returned nil")
			}

			back := ParseGrpcError(grpcErr)
			if back == nil {
				t.Fatal("ParseGrpcError returned nil")
			}

			if back.Code != tc.code {
				t.Fatalf("round trip: got HTTP %d, want %d", back.Code, tc.code)
			}
		})
	}
}

func TestToGrpcErrorCodeMapping(t *testing.T) {
	cases := []struct {
		httpCode int
		wantCode codes.Code
	}{
		{http.StatusTooManyRequests, codes.ResourceExhausted},
		{http.StatusServiceUnavailable, codes.Unavailable},
		{http.StatusGatewayTimeout, codes.DeadlineExceeded},
		{http.StatusConflict, codes.AlreadyExists},
		{http.StatusForbidden, codes.PermissionDenied},
	}

	for _, tc := range cases {
		err := ToGrpcError(&AppError{Code: tc.httpCode, Message: "x"})
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("HTTP %d: not a gRPC status error", tc.httpCode)
		}
		if st.Code() != tc.wantCode {
			t.Fatalf("HTTP %d -> got %s, want %s", tc.httpCode, st.Code(), tc.wantCode)
		}
	}
}

func TestParseGrpcErrorFallbackMapping(t *testing.T) {
	cases := []struct {
		grpcCode codes.Code
		wantCode int
	}{
		{codes.ResourceExhausted, http.StatusTooManyRequests},
		{codes.Unavailable, http.StatusServiceUnavailable},
		{codes.DeadlineExceeded, http.StatusGatewayTimeout},
		{codes.AlreadyExists, http.StatusConflict},
		{codes.Unauthenticated, http.StatusUnauthorized},
	}

	for _, tc := range cases {
		err := status.Error(tc.grpcCode, "x")
		appErr := ParseGrpcError(err)
		if appErr == nil {
			t.Fatalf("%s: ParseGrpcError returned nil", tc.grpcCode)
		}
		if appErr.Code != tc.wantCode {
			t.Fatalf("%s -> got HTTP %d, want %d", tc.grpcCode, appErr.Code, tc.wantCode)
		}
	}
}

func TestInvalidAccessTokenIsUnauthorized(t *testing.T) {
	err := InvalidAccessToken()
	var apiErr *AppError
	if !errorsAs(err, &apiErr) {
		t.Fatalf("InvalidAccessToken() should return an AppError, got %T", err)
	}
	if apiErr.Code != http.StatusUnauthorized {
		t.Fatalf("got HTTP %d, want 401", apiErr.Code)
	}
}

// errorsAs is a tiny alias to avoid importing errors twice with a name clash.
func errorsAs(err error, target **AppError) bool {
	apiErr, ok := err.(*AppError)
	if ok {
		*target = apiErr
	}
	return ok
}
