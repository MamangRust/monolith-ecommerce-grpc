package banner_test

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

// gapi: non-existent banner must map to codes.NotFound (404), not Internal.
func (s *BannerGapiTestSuite) TestBannerGapiNotFound() {
	ctx := context.Background()
	_, err := s.queryClient.FindById(ctx, &pb.FindByIdBannerRequest{Id: 999999})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "non-existent banner must be NotFound, got %v: %s", st.Code(), st.Message())
}

// api: non-existent banner must map to 404, invalid path ID to 400.
func (s *BannerApiTestSuite) TestBannerApiNotFound() {
	apigatewaymiddlewares.RegisterErrorHandler(s.echo)

	req := httptest.NewRequest(http.MethodGet, "/api/banner-query/999999", nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusNotFound, rec.Code, "non-existent banner must be 404, got %d: %s", rec.Code, rec.Body.String())
}

func (s *BannerApiTestSuite) TestBannerApiInvalidID() {
	apigatewaymiddlewares.RegisterErrorHandler(s.echo)

	req := httptest.NewRequest(http.MethodGet, "/api/banner-query/abc", nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code, "invalid banner ID must be 400, got %d: %s", rec.Code, rec.Body.String())
}

// repository: FindByID on a non-existent ID must return a typed not-found error.
func (s *BannerRepositoryTestSuite) TestBannerFindByIDNotFound() {
	ctx := context.Background()
	_, err := s.repo.BannerQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}
