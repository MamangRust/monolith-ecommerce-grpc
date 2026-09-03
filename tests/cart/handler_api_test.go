package cart_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	carthandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/cart"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type CartApiTestSuite struct {
	tests.BaseTestSuite
	echo    *echo.Echo
	userID  int
	cartIDs []int
}

func (s *CartApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Bootstrap all required downstream services
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupCartService()

	// Seed a user to use as the authenticated caller
	ctx := s.Ctx
	s.userID = s.SeedUser(ctx)

	// Build Echo router with an auth middleware that injects the seeded user_id
	s.echo = echo.New()
	s.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user_id", s.userID)
			return next(c)
		}
	})

	// Cache infrastructure used by the cart API handler
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)

	// Register the cart HTTP routes (they connect to the cart gRPC service)
	carthandler.RegisterCartHandler(&carthandler.DepsCart{
		Client:     s.Conns["cart"],
		E:          s.echo,
		Logger:     s.Log,
		CacheStore: cacheStore,
	})

	s.cartIDs = make([]int, 0)
}

func (s *CartApiTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *CartApiTestSuite) TestCartApiLifecycle() {
	ctx := s.Ctx

	// Seed entities needed for cart items
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, s.userID)
	prod1ID := s.SeedProduct(ctx, merchantID, categoryID)
	prod2ID := s.SeedProduct(ctx, merchantID, categoryID)

	// ────────────────────────────────────────────
	// 1. Create — add first product to cart
	// ────────────────────────────────────────────
	createBody := requests.CreateCartRequest{
		ProductID: prod1ID,
		Quantity:  2,
	}
	body, _ := json.Marshal(createBody)
	req := httptest.NewRequest(http.MethodPost, "/api/cart-command/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusCreated, rec.Code, "create first cart item: %s", rec.Body.String())

	var createRes map[string]interface{}
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &createRes))
	data := createRes["data"].(map[string]interface{})
	cart1ID := int(data["id"].(float64))
	s.cartIDs = append(s.cartIDs, cart1ID)
	s.Equal(float64(prod1ID), data["product_id"])
	s.Equal(float64(2), data["quantity"])

	// ────────────────────────────────────────────
	// 2. Create — add second product to cart
	// ────────────────────────────────────────────
	createBody2 := requests.CreateCartRequest{
		ProductID: prod2ID,
		Quantity:  1,
	}
	body2, _ := json.Marshal(createBody2)
	req2 := httptest.NewRequest(http.MethodPost, "/api/cart-command/create", bytes.NewBuffer(body2))
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec2 := httptest.NewRecorder()
	s.echo.ServeHTTP(rec2, req2)
	s.Require().Equal(http.StatusCreated, rec2.Code, "create second cart item: %s", rec2.Body.String())

	var createRes2 map[string]interface{}
	s.Require().NoError(json.Unmarshal(rec2.Body.Bytes(), &createRes2))
	data2 := createRes2["data"].(map[string]interface{})
	cart2ID := int(data2["id"].(float64))
	s.cartIDs = append(s.cartIDs, cart2ID)

	// ────────────────────────────────────────────
	// 3. FindAll — list cart items for the user
	// ────────────────────────────────────────────
	req3 := httptest.NewRequest(http.MethodGet, "/api/cart-query?page=1&page_size=10", nil)
	rec3 := httptest.NewRecorder()
	s.echo.ServeHTTP(rec3, req3)
	s.Require().Equal(http.StatusOK, rec3.Code, "find all cart items: %s", rec3.Body.String())

	var findAllRes map[string]interface{}
	s.Require().NoError(json.Unmarshal(rec3.Body.Bytes(), &findAllRes))
	findAllData := findAllRes["data"].([]interface{})
	s.GreaterOrEqual(len(findAllData), 2, "should have at least 2 cart items")

	// ────────────────────────────────────────────
	// 4. Delete — remove the first cart item
	// ────────────────────────────────────────────
	deleteBody := requests.DeleteCartRequest{
		CartID: cart1ID,
		UserID: s.userID,
	}
	body4, _ := json.Marshal(deleteBody)
	req4 := httptest.NewRequest(http.MethodDelete, "/api/cart-command/delete", bytes.NewBuffer(body4))
	req4.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec4 := httptest.NewRecorder()
	s.echo.ServeHTTP(rec4, req4)
	s.Require().Equal(http.StatusOK, rec4.Code, "delete cart item: %s", rec4.Body.String())
	s.cartIDs = s.cartIDs[1:] // remove from tracking

	// ────────────────────────────────────────────
	// 5. DeleteAll — remove remaining cart items
	// ────────────────────────────────────────────
	deleteAllBody := requests.DeleteAllCartRequest{
		UserID:  s.userID,
		CartIds: s.cartIDs,
	}
	body5, _ := json.Marshal(deleteAllBody)
	req5 := httptest.NewRequest(http.MethodPost, "/api/cart-command/delete-all", bytes.NewBuffer(body5))
	req5.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec5 := httptest.NewRecorder()
	s.echo.ServeHTTP(rec5, req5)
	s.Require().Equal(http.StatusOK, rec5.Code, "delete all cart items: %s", rec5.Body.String())

	// ────────────────────────────────────────────
	// 6. Verify — FindAll should now be empty
	// ────────────────────────────────────────────
	req6 := httptest.NewRequest(http.MethodGet, "/api/cart-query?page=1&page_size=10", nil)
	rec6 := httptest.NewRecorder()
	s.echo.ServeHTTP(rec6, req6)
	s.Require().Equal(http.StatusOK, rec6.Code)

	var verifyRes map[string]interface{}
	s.Require().NoError(json.Unmarshal(rec6.Body.Bytes(), &verifyRes))
	verifyData, _ := verifyRes["data"].([]interface{}) // nil (empty) and [] both mean empty
	s.Empty(verifyData, "cart should be empty after deleting all items")
}

func TestCartApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartApiTestSuite))
}
