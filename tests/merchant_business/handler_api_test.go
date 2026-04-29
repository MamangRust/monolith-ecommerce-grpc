package merchant_business_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	merchantbusinesshandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/merchant_business"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type MerchantBusinessApiTestSuite struct {
	tests.BaseTestSuite
	echo       *echo.Echo
	businessID int
	merchantID int
	userID     int
}

func (s *MerchantBusinessApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupMerchantBusinessService()

	// Seed dependencies
	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)

	s.echo = echo.New()

	merchantbusinesshandler.RegisterMerchantBusinessHandler(&merchantbusinesshandler.DepsMerchantBusiness{
		Client:     s.Conns["merchant_business"],
		E:          s.echo,
		Logger:     s.Log,
		CacheStore: s.GetCacheStore(),
	})
}

func (s *MerchantBusinessApiTestSuite) TestMerchantBusinessApiLifecycle() {
	// 1. Create
	reqBody := requests.CreateMerchantBusinessInformationRequest{
		MerchantID:        s.merchantID,
		BusinessType:      "Retail",
		TaxID:             "123-456-789",
		EstablishedYear:   2020,
		NumberOfEmployees: 10,
		WebsiteUrl:        "http://example.com",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/merchant-business-command/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal(reqBody.BusinessType, data["business_type"])
	s.businessID = int(data["id"].(float64))

	// 2. FindById
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant-business-query/%d", s.businessID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.businessID), data["id"])

	// 3. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/merchant-business-query", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/merchant-business-query/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Update
	updateBody := requests.UpdateMerchantBusinessInformationRequest{
		BusinessType:      "Wholesale",
		TaxID:             "987-654-321",
		EstablishedYear:   2021,
		NumberOfEmployees: 20,
		WebsiteUrl:        "http://updated.com",
	}
	body, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-business-command/update/%d", s.businessID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(updateBody.BusinessType, data["business_type"])

	// 6. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-business-command/trashed/%d", s.businessID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/merchant-business-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-business-command/restore/%d", s.businessID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-business-command/trashed/%d", s.businessID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/merchant-business-command/permanent/%d", s.businessID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/merchant-business-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/merchant-business-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestMerchantBusinessApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantBusinessApiTestSuite))
}
