package review_detail_test

import (
	"context"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"

	reviewdetailhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/review_detail"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type ReviewDetailApiTestSuite struct {
	tests.BaseTestSuite
	echo           *echo.Echo
	reviewDetailID int
	reviewID       int
}

func (s *ReviewDetailApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupOrderService()
	s.SetupReviewService()
	s.SetupReviewDetailService()

	s.echo = echo.New()

	// Seed dependencies
	ctx := context.Background()
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)
	// SeedOrder(ctx context.Context, userID int, merchID int, prodID int) int
	s.SeedOrder(ctx, userID, merchID, prodID)
	s.reviewID = s.SeedReview(ctx, userID, prodID)

	reviewdetailhandler.RegisterReviewDetailHandler(&reviewdetailhandler.DepsReviewDetail{
		Client:        s.Conns["review-detail"],
		E:             s.echo,
		Logger:        s.Log,
		Cache:         s.GetCacheStore(),
		Upload:        &tests.MockImageUpload{},
		Observability: s.Obs,
	})
}

func (s *ReviewDetailApiTestSuite) Test01_CreateReviewDetail() {
	fields := map[string]string{
		"review_id": fmt.Sprintf("%d", s.reviewID),
		"type":      "photo",
		"caption":   "API Test Caption",
	}
	body, contentType := s.BuildMultipartRequestBody(fields, "url", "review_photo.jpg")

	req := httptest.NewRequest(http.MethodPost, "/api/review-detail-command/create", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusCreated, rec.Code, rec.Body.String())

	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal("API Test Caption", data["caption"])
	s.reviewDetailID = int(data["id"].(float64))
}

func (s *ReviewDetailApiTestSuite) Test02_FindById() {
	s.Require().NotZero(s.reviewDetailID)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/review-detail-query/%d", s.reviewDetailID), nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)

	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal(float64(s.reviewDetailID), data["id"])
}

func (s *ReviewDetailApiTestSuite) Test03_FindAll() {
	req := httptest.NewRequest(http.MethodGet, "/api/review-detail-query", nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *ReviewDetailApiTestSuite) Test04_FindByActive() {
	req := httptest.NewRequest(http.MethodGet, "/api/review-detail-query/active", nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *ReviewDetailApiTestSuite) Test05_UpdateReviewDetail() {
	s.Require().NotZero(s.reviewDetailID)
	fields := map[string]string{
		"type":    "video",
		"caption": "Updated Caption",
	}
	body, contentType := s.BuildMultipartRequestBody(fields, "url", "review_video.mp4")

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review-detail-command/update/%d", s.reviewDetailID), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *ReviewDetailApiTestSuite) Test06_TrashedReviewDetail() {
	s.Require().NotZero(s.reviewDetailID)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review-detail-command/trashed/%d", s.reviewDetailID), nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *ReviewDetailApiTestSuite) Test07_FindByTrashed() {
	req := httptest.NewRequest(http.MethodGet, "/api/review-detail-query/trashed", nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *ReviewDetailApiTestSuite) Test08_RestoreReviewDetail() {
	s.Require().NotZero(s.reviewDetailID)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review-detail-command/restore/%d", s.reviewDetailID), nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *ReviewDetailApiTestSuite) Test09_DeletePermanent() {
	s.Require().NotZero(s.reviewDetailID)
	
	// Must be trashed before permanent deletion according to service logic
	trashReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/review-detail-command/trashed/%d", s.reviewDetailID), nil)
	trashRec := httptest.NewRecorder()
	s.echo.ServeHTTP(trashRec, trashReq)
	s.Require().Equal(http.StatusOK, trashRec.Code)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/review-detail-command/permanent/%d", s.reviewDetailID), nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *ReviewDetailApiTestSuite) Test10_RestoreAll() {
	req := httptest.NewRequest(http.MethodPost, "/api/review-detail-command/restore/all", nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *ReviewDetailApiTestSuite) Test11_DeleteAll() {
	req := httptest.NewRequest(http.MethodPost, "/api/review-detail-command/permanent/all", nil)
	rec := httptest.NewRecorder()

	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func TestReviewDetailApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewDetailApiTestSuite))
}
