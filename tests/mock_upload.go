package tests

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type MockImageUpload struct{}

func (m *MockImageUpload) EnsureUploadDirectory(uploadDir string) error {
	return nil
}

func (m *MockImageUpload) ProcessImageUpload(c echo.Context, uploadDir string, file *multipart.FileHeader, isDocument bool) (string, error) {
	return "http://example.com/mock-image.jpg", nil
}

func (m *MockImageUpload) CleanupImageOnFailure(imagePath string) {}

func (m *MockImageUpload) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	return nil
}
