package cart_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	apigatewaymiddlewares "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/middlewares"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gapi: cart has no FindById query; the valid negative path is validation of
// create payload — quantity 0 must be rejected as InvalidArgument.
func (s *CartGapiTestSuite) TestCartGapiInvalidQuantity() {
	ctx := context.Background()
	_, err := s.commandClient.Create(ctx, &pb.CreateCartRequest{
		UserId:    1,
		ProductId: 1,
		Quantity:  0,
	})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.InvalidArgument, st.Code(), "cart create with quantity 0 must be InvalidArgument, got %v: %s", st.Code(), st.Message())
}

// api: malformed JSON body on cart create must map to 400.
func (s *CartApiTestSuite) TestCartApiInvalidBody() {
	apigatewaymiddlewares.RegisterErrorHandler(s.echo)

	req := httptest.NewRequest(http.MethodPost, "/api/cart-command/create", strings.NewReader("{not-json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code, "malformed cart create body must be 400, got %d: %s", rec.Code, rec.Body.String())
}
