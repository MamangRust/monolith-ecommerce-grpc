package order_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	orderhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/order"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type OrderApiTestSuite struct {
	tests.BaseTestSuite
	echo    *echo.Echo
	orderID int
	userID  int
	merchID int
	prodID  int
}

func (s *OrderApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupCategoryService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()
	s.SetupOrderService()

	// Seed dependencies
	ctx := context.Background()
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	s.userID = userID
	s.merchID = merchID
	s.prodID = prodID

	s.echo = echo.New()

	orderhandler.RegisterOrderHandler(&orderhandler.DepsOrder{
		Client:     s.Conns["order"],
		E:          s.echo,
		Logger:     s.Log,
		CacheStore: s.GetCacheStore(),
	})
}

func (s *OrderApiTestSuite) TestOrderApiLifecycle() {
	// 1. Create
	reqBody := requests.CreateOrderRequest{
		MerchantID: s.merchID,
		UserID:     s.userID,
		TotalPrice: 1000,
		Items: []requests.CreateOrderItemRequest{
			{ProductID: s.prodID, Quantity: 1, Price: 1000},
		},
		ShippingAddress: requests.CreateShippingAddressRequest{
			Alamat:         "Test Alamat",
			Provinsi:       "Test Provinsi",
			Kota:           "Test Kota",
			Courier:        "Test Courier",
			ShippingMethod: "Test Method",
			ShippingCost:   100,
			Negara:         "Test Negara",
		},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/order-command/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.orderID = int(data["id"].(float64))

	// 2. FindById
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/order-query/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.orderID), data["id"])

	// 3. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/order-query", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/order-query/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Update
	updateBody := requests.UpdateOrderRequest{
		UserID:     s.userID,
		TotalPrice: 1500,
		Items: []requests.UpdateOrderItemRequest{
			{ProductID: s.prodID, Quantity: 1, Price: 1500},
		},
		ShippingAddress: &requests.UpdateShippingAddressRequest{
			Alamat:         "Updated Alamat",
			Provinsi:       "Updated Provinsi",
			Kota:           "Updated Kota",
			Courier:        "Updated Courier",
			ShippingMethod: "Updated Method",
			ShippingCost:   200,
			Negara:         "Updated Negara",
		},
	}
	body, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order-command/update/%d", s.orderID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order-command/trashed/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/order-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order-command/restore/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order-command/trashed/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/order-command/permanent/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/order-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/order-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *OrderApiTestSuite) Test12_GetMonthlyTotalRevenue() {
	req := httptest.NewRequest(http.MethodGet, "/api/order/monthly-total-revenue?year=2024&month=1", nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *OrderApiTestSuite) Test13_GetYearlyTotalRevenue() {
	req := httptest.NewRequest(http.MethodGet, "/api/order/yearly-total-revenue?year=2024", nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func TestOrderApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderApiTestSuite))
}
