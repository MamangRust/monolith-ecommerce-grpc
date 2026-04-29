package transaction_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"bytes"
	"fmt"
	"net/url"
	"strings"

	transactionhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/transaction"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type TransactionApiTestSuite struct {
	tests.BaseTestSuite
	echo          *echo.Echo
	transactionID int
	userID        int
	merchID       int
	orderID       int
}

func (s *TransactionApiTestSuite) SetupSuite() {
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
	merchID := s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, merchID, catID)
	orderID := s.SeedOrder(ctx, userID, merchID, prodID)
	s.SeedShippingAddress(ctx, orderID)
	s.SeedOrderItem(ctx, orderID, prodID)

	s.userID = userID
	s.merchID = merchID
	s.orderID = orderID

	s.echo = echo.New()

	transactionhandler.RegisterTransactionHandler(&transactionhandler.DepsTransaction{
		Client:     s.Conns["transaction"],
		E:          s.echo,
		Logger:     s.Log,
		CacheStore: s.GetCacheStore(),
		ApiHandler: errors.NewApiHandler(s.Obs, s.Log),
	})
}

func (s *TransactionApiTestSuite) TestTransactionApiLifecycle() {
	// 1. Create
	fields := map[string]interface{}{
		"order_id":    s.orderID,
		"merchant_id": s.merchID,
		"amount":      100000,
		"user_id":     s.userID,
		"status":      "pending",
		"method":      "credit_card",
	}
	body, _ := json.Marshal(fields)
	req := httptest.NewRequest(http.MethodPost, "/api/transaction-command/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusCreated, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.transactionID = int(data["id"].(float64))

	// 2. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/transaction-query", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. FindById
	s.Require().NotZero(s.transactionID)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/transaction-query/%d", s.transactionID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.transactionID), data["id"])

	// 4. FindByMerchant
	s.Require().NotZero(s.merchID)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/transaction-query/merchant/%d", s.merchID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/transaction-query/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Update
	form := url.Values{}
	form.Add("status", "success")
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transaction-command/update/%d", s.transactionID), strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transaction-command/trashed/%d", s.transactionID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/transaction-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transaction-command/restore/%d", s.transactionID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transaction-command/trashed/%d", s.transactionID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/transaction-command/permanent/%d", s.transactionID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/transaction-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 12. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/transaction-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *TransactionApiTestSuite) Test13_MonthlySuccessStats() {
	req := httptest.NewRequest(http.MethodGet, "/api/transaction-stats/monthly-success?year=2026&month=4", nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func TestTransactionApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionApiTestSuite))
}
