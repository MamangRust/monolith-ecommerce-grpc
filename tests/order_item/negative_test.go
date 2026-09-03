package order_item_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	apigatewaymiddlewares "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/middlewares"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

// order_item has no single-record lookup: FindOrderItemByOrder returns an empty
// list (success) for a non-existent order, so there is no NotFound path.
//
// gapi: a non-existent order must return an empty result, not an error.
func (s *OrderItemGapiTestSuite) TestOrderItemGapiEmptyResult() {
	ctx := context.Background()
	res, err := s.queryClient.FindOrderItemByOrder(ctx, &pb.FindByIdOrderItemRequest{Id: 999999})
	s.NoError(err)
	s.NotNil(res)
	s.Empty(res.Data)
}

// api: a non-existent order returns 200 with empty data; invalid ID maps to 400.
func (s *OrderItemApiTestSuite) TestOrderItemApiEmptyResult() {
	apigatewaymiddlewares.RegisterErrorHandler(s.echo)

	req := httptest.NewRequest(http.MethodGet, "/api/order-item/order/999999", nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code, "non-existent order must return 200 with empty data, got %d: %s", rec.Code, rec.Body.String())
}

func (s *OrderItemApiTestSuite) TestOrderItemApiInvalidID() {
	apigatewaymiddlewares.RegisterErrorHandler(s.echo)

	req := httptest.NewRequest(http.MethodGet, "/api/order-item/order/abc", nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code, "invalid order ID must be 400, got %d: %s", rec.Code, rec.Body.String())
}

// repository: FindOrderItemByOrder on a non-existent order must return an empty
// result without error.
func (s *OrderItemRepositoryTestSuite) TestOrderItemFindByOrderEmpty() {
	ctx := context.Background()
	items, err := s.repo.OrderItemQuery.FindOrderItemByOrder(ctx, 999999)
	s.NoError(err)
	s.Empty(items)
}
