package order_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	apigatewaymiddlewares "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/middlewares"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gapi: non-existent order must map to codes.NotFound (404), not Internal.
func (s *OrderGapiTestSuite) TestOrderGapiNotFound() {
	ctx := context.Background()
	_, err := s.queryClient.FindById(ctx, &pb.FindByIdOrderRequest{Id: 999999})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "non-existent order must be NotFound, got %v: %s", st.Code(), st.Message())
}

// api: non-existent order must map to 404, invalid path ID to 400.
func (s *OrderApiTestSuite) TestOrderApiNotFound() {
	apigatewaymiddlewares.RegisterErrorHandler(s.echo)

	req := httptest.NewRequest(http.MethodGet, "/api/order-query/999999", nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusNotFound, rec.Code, "non-existent order must be 404, got %d: %s", rec.Code, rec.Body.String())
}

func (s *OrderApiTestSuite) TestOrderApiInvalidID() {
	apigatewaymiddlewares.RegisterErrorHandler(s.echo)

	req := httptest.NewRequest(http.MethodGet, "/api/order-query/abc", nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code, "invalid order ID must be 400, got %d: %s", rec.Code, rec.Body.String())
}

// repository: FindByID on a non-existent ID must return a typed not-found error.
func (s *OrderRepositoryTestSuite) TestOrderFindByIDNotFound() {
	ctx := context.Background()
	_, err := s.repo.OrderQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}
