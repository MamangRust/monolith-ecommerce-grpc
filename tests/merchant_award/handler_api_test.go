package merchant_award_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	merchantawardhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/merchant_award"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type MerchantAwardApiTestSuite struct {
	tests.BaseTestSuite
	echo       *echo.Echo
	awardID    int
	merchantID int
}

func (s *MerchantAwardApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupMerchantAwardService()

	// Seed dependencies
	ctx := context.Background()
	userID := s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, userID)

	s.echo = echo.New()

	merchantawardhandler.RegisterMerchantAwardHandler(&merchantawardhandler.DepsMerchantAward{
		Client:     s.Conns["merchant_award"],
		E:          s.echo,
		Logger:     s.Log,
		CacheStore: s.GetCacheStore(),
	})
}

func (s *MerchantAwardApiTestSuite) TestMerchantAwardApiLifecycle() {
	// 1. Create
	reqBody := requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID:     s.merchantID,
		Title:          "Best Merchant 2024",
		Description:    "Award for excellence",
		IssuedBy:       "E-commerce Platform",
		IssueDate:      "2024-01-01",
		ExpiryDate:     "2025-01-01",
		CertificateUrl: "http://example.com/cert.pdf",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/merchant-award-command/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal(reqBody.Title, data["title"])
	s.awardID = int(data["id"].(float64))

	// 2. FindById
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant-award-query/%d", s.awardID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.awardID), data["id"])

	// 3. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/merchant-award-query", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/merchant-award-query/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Update
	updateBody := requests.UpdateMerchantCertificationOrAwardRequest{
		Title:          "Updated Award Title",
		Description:    "Updated Description",
		IssuedBy:       "Updated Issuer",
		IssueDate:      "2024-02-01",
		ExpiryDate:     "2025-02-01",
		CertificateUrl: "http://example.com/updated.pdf",
	}
	body, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-award-command/update/%d", s.awardID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(updateBody.Title, data["title"])

	// 6. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-award-command/trashed/%d", s.awardID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/merchant-award-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-award-command/restore/%d", s.awardID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-award-command/trashed/%d", s.awardID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/merchant-award-command/permanent/%d", s.awardID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/merchant-award-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/merchant-award-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestMerchantAwardApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantAwardApiTestSuite))
}
