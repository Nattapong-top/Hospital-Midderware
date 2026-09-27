package testutil

import (
	"net/http"
	"net/http/httptest"

	"Hospital-Middleware/internal/domain"

	"github.com/gin-gonic/gin"
)

// SetupGinTestContext creates a test Gin context and recorder for HTTP handler tests.
func SetupGinTestContext(req *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	return c, rec
}

// MockAdapter is a reusable mock for domain.ExternalAPIAdapter.
type MockAdapter struct {
	Patient *domain.PatientDTO
	Err     error
}

func (m *MockAdapter) Search(criteria domain.SearchCriteria) (*domain.PatientDTO, error) {
	return m.Patient, m.Err
}

// MockResolver is a reusable mock for domain.HospitalAPIResolver.
type MockResolver struct {
	Adapter domain.ExternalAPIAdapter
	Err     error
}

func (m *MockResolver) Resolve(hospitalID string) (domain.ExternalAPIAdapter, error) {
	return m.Adapter, m.Err
}
