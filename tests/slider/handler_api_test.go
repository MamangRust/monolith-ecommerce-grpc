package slider_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sliderhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/slider"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type SliderApiTestSuite struct {
	tests.BaseTestSuite
	echo    *echo.Echo
	sliderID int
}

func (s *SliderApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupSliderService()

	s.echo = echo.New()

	sliderhandler.RegisterSliderHandler(&sliderhandler.DepsSlider{
		Client: s.Conns["slider"],
		E:      s.echo,
		Logger: s.Log,
		Cache:  s.GetCacheStore(),
		Upload: &tests.MockImageUpload{},
	})
}

func (s *SliderApiTestSuite) TestSliderApiLifecycle() {
	// 1. Create
	fields := map[string]string{
		"name": "Test Slider",
	}
	body, contentType := s.BuildMultipartRequestBody(fields, "image_slider", "slider.jpg")
	req := httptest.NewRequest(http.MethodPost, "/api/slider-command/create", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusCreated, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal("Test Slider", data["name"])
	s.sliderID = int(data["id"].(float64))

	// 2. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/slider-query", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/slider-query/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Update
	s.Require().NotZero(s.sliderID)
	updateFields := map[string]string{
		"name": "Updated Test Slider",
	}
	body, contentType = s.BuildMultipartRequestBody(updateFields, "image_slider", "slider_updated.jpg")
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/slider-command/update/%d", s.sliderID), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, contentType)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/slider-command/trashed/%d", s.sliderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/slider-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/slider-command/restore/%d", s.sliderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/slider-command/trashed/%d", s.sliderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/slider-command/permanent/%d", s.sliderID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/slider-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/slider-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestSliderApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SliderApiTestSuite))
}
