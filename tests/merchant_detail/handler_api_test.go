package merchant_detail_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	merchantdetailhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/merchant_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type MerchantDetailApiTestSuite struct {
	tests.BaseTestSuite
	echo       *echo.Echo
	detailID   int
	merchantID int
	userID     int
}

func (s *MerchantDetailApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupMerchantDetailService()

	// Seed dependencies
	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)

	s.echo = echo.New()
	apiHandler := errors.NewApiHandler(s.Obs, s.Log)

	merchantdetailhandler.RegisterMerchantDetailHandler(&merchantdetailhandler.DepsMerchantDetail{
		Client:      s.Conns["merchant_detail"],
		E:           s.echo,
		Logger:      s.Log,
		CacheStore:  s.GetCacheStore(),
		UploadImage: &tests.MockImageUpload{},
		ApiHandler:  apiHandler,
	})
}

func (s *MerchantDetailApiTestSuite) TestMerchantDetailApiLifecycle() {
	// 1. Create
	reqBody := requests.CreateMerchantDetailRequest{
		MerchantID:       s.merchantID,
		DisplayName:      "Test Merchant",
		CoverImageUrl:    "http://example.com/cover.jpg",
		LogoUrl:          "http://example.com/logo.jpg",
		ShortDescription: "A test merchant",
		WebsiteUrl:       "http://example.com",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/merchant-detail-command/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusCreated, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal(reqBody.DisplayName, data["display_name"])
	s.detailID = int(data["id"].(float64))

	// 2. FindById
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant-detail-query/%d", s.detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.detailID), data["id"])

	// 3. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/merchant-detail-query", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/merchant-detail-query/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Update
	updateBody := requests.UpdateMerchantDetailRequest{
		DisplayName:      "Updated Merchant",
		CoverImageUrl:    "http://example.com/updated_cover.jpg",
		LogoUrl:          "http://example.com/updated_logo.jpg",
		ShortDescription: "Updated short description",
		WebsiteUrl:       "http://updated.com",
	}
	body, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-detail-command/update/%d", s.detailID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(updateBody.DisplayName, data["display_name"])

	// 6. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-detail-command/trashed/%d", s.detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/merchant-detail-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-detail-command/restore/%d", s.detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeletePermanent (Trash first)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-detail-command/trashed/%d", s.detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/merchant-detail-command/permanent/%d", s.detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().True(rec.Code == http.StatusOK || rec.Code == http.StatusNoContent)

	// 10. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/merchant-detail-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/merchant-detail-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestMerchantDetailApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDetailApiTestSuite))
}
