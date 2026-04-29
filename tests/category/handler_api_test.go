package category_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	categoryhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/category"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type CategoryApiTestSuite struct {
	tests.BaseTestSuite
	echo       *echo.Echo
	categoryID int
}

func (s *CategoryApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupCategoryService()

	s.echo = echo.New()
	apiHandler := errors.NewApiHandler(s.Obs, s.Log)

	categoryhandler.RegisterCategoryHandler(&categoryhandler.DepsCategory{
		Client:      s.Conns["category"],
		E:           s.echo,
		Logger:      s.Log,
		CacheStore:  s.GetCacheStore(),
		UploadImage: &tests.MockImageUpload{},
		ApiHandler:  apiHandler,
	})
}

func (s *CategoryApiTestSuite) TestCategoryApiLifecycle() {
	// 1. Create
	fields := map[string]string{
		"name":           "Test Category",
		"description":    "Test Description",
		"slug_category":  "test-category",
		"image_category": "test.jpg",
	}
	body, contentType := s.BuildMultipartRequestBody(fields, "image", "test.jpg")
	req := httptest.NewRequest(http.MethodPost, "/api/category-command/create", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusCreated, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal(fields["name"], data["name"])
	s.categoryID = int(data["id"].(float64))

	// 2. FindById
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/category-query/%d", s.categoryID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.categoryID), data["id"])

	// 3. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/category-query", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/category-query/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Update
	updFields := map[string]string{
		"name":           "Updated Category",
		"description":    "Updated Description",
		"slug_category":  "updated-category",
		"image_category": "updated.jpg",
	}
	updBody, updContentType := s.BuildMultipartRequestBody(updFields, "image", "updated.jpg")
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/category-command/update/%d", s.categoryID), bytes.NewReader(updBody))
	req.Header.Set(echo.HeaderContentType, updContentType)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code, rec.Body.String())
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(updFields["name"], data["name"])

	// 6. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/category-command/trashed/%d", s.categoryID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/category-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/category-command/restore/%d", s.categoryID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/category-command/trashed/%d", s.categoryID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/category-command/permanent/%d", s.categoryID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/category-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. DeleteAll
	req = httptest.NewRequest(http.MethodDelete, "/api/category-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestCategoryApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryApiTestSuite))
}
