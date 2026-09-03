package review_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	apigatewaymiddlewares "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/middlewares"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gapi: review has no FindById query RPC; Update on a non-existent review must
// map to codes.NotFound (404), not Internal.
func (s *ReviewGapiTestSuite) TestReviewGapiNotFound() {
	ctx := context.Background()
	_, err := s.commandClient.Update(ctx, &pb.UpdateReviewRequest{
		ReviewId: 999999,
		Rating:   5,
		Comment:  "not found",
	})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "update on non-existent review must be NotFound, got %v: %s", st.Code(), st.Message())
}

// api: update on a non-existent review must map to 404, invalid path ID to 400.
func (s *ReviewApiTestSuite) TestReviewApiNotFound() {
	apigatewaymiddlewares.RegisterErrorHandler(s.echo)

	// Valid body required — otherwise gateway validation (400) fires before the NotFound lookup.
	body := strings.NewReader(`{"rating": 5, "comment": "not found"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/review-command/update/999999", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusNotFound, rec.Code, "update on non-existent review must be 404, got %d: %s", rec.Code, rec.Body.String())
}

func (s *ReviewApiTestSuite) TestReviewApiInvalidID() {
	apigatewaymiddlewares.RegisterErrorHandler(s.echo)

	req := httptest.NewRequest(http.MethodPost, "/api/review-command/update/abc", nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code, "invalid review ID must be 400, got %d: %s", rec.Code, rec.Body.String())
}

// repository: FindByID on a non-existent ID must return a typed not-found error.
func (s *ReviewRepositoryTestSuite) TestReviewFindByIDNotFound() {
	ctx := context.Background()
	_, err := s.repo.ReviewQuery.FindByID(ctx, 999999)
	s.Require().Error(err)
	var appErr *errors.AppError
	s.Require().ErrorAs(err, &appErr)
	s.Equal(errors.ErrorTypeNotFound, appErr.Type, "expected not-found error type, got %s: %v", appErr.Type, err)
}
