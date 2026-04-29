package order_item_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	orderitemhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/order_item"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type OrderItemApiTestSuite struct {
	tests.BaseTestSuite
	echo    *echo.Echo
	orderID int
}

func (s *OrderItemApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupShippingAddressService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()

	// Seed dependencies
	ctx := context.Background()
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)
	s.orderID = s.SeedOrder(ctx, userID, merchID, prodID)

	s.echo = echo.New()

	orderitemhandler.RegisterOrderItemHandler(&orderitemhandler.DepsOrderItem{
		Client:     s.Conns["order-item"],
		E:          s.echo,
		Logger:     s.Log,
		CacheStore: s.GetCacheStore(),
	})
}

func (s *OrderItemApiTestSuite) TestOrderItemApiLifecycle() {
	// 1. FindAll
	req := httptest.NewRequest(http.MethodGet, "/api/order-item", nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 2. FindOrderItemByOrder
	s.Require().NotZero(s.orderID)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/order-item/order/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].([]interface{})
	if len(data) == 0 {
		s.T().Skip("No order items found")
	}
	itemID := int(data[0].(map[string]interface{})["id"].(float64))

	// 3. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/order-item/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/order-item/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order-item/trash/%d", itemID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order-item/restore/%d", itemID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order-item/trash/%d", itemID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/order-item/permanent/%d", itemID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/order-item/restore/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/order-item/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestOrderItemApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemApiTestSuite))
}
