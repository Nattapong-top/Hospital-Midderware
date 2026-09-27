package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"Hospital-Middleware/internal/application"
	httpDelivery "Hospital-Middleware/internal/delivery/http"
	"Hospital-Middleware/internal/domain"
	"Hospital-Middleware/internal/testutil"

	"github.com/gin-gonic/gin"
)

type mockAdapter = testutil.MockAdapter
type mockResolver = testutil.MockResolver

func TestPatientHandler_SearchPatient_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPat := &domain.PatientDTO{
		FirstNameTH: "สมชาย",
		NationalID:  "1100200300400",
		PatientHN:   "HN-12345",
	}
	adapter := &mockAdapter{Patient: mockPat}
	resolver := &mockResolver{Adapter: adapter}
	searchUC := application.NewSearchPatient(resolver)
	handler := httpDelivery.NewPatientHandler(searchUC)

	r := gin.New()
	r.GET("/patient/search", func(c *gin.Context) {
		c.Set("hospital", "HN99999")
		handler.SearchPatient(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/patient/search?national_id=1100200300400", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("คาดหวัง Status 200 แต่ได้ %d (Response: %s)", rec.Code, rec.Body.String())
	}
}

func TestPatientHandler_SearchPatient_MissingHospitalContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resolver := &mockResolver{}
	searchUC := application.NewSearchPatient(resolver)
	handler := httpDelivery.NewPatientHandler(searchUC)

	r := gin.New()
	r.GET("/patient/search", handler.SearchPatient)

	req := httptest.NewRequest(http.MethodGet, "/patient/search?national_id=1100200300400", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("คาดหวัง Status 401 แต่ได้ %d", rec.Code)
	}
}

func TestPatientHandler_SearchPatient_InvalidCriteria(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resolver := &mockResolver{}
	searchUC := application.NewSearchPatient(resolver)
	handler := httpDelivery.NewPatientHandler(searchUC)

	r := gin.New()
	r.GET("/patient/search", func(c *gin.Context) {
		c.Set("hospital", "HN99999")
		handler.SearchPatient(c)
	})

	// ส่ง request แบบไม่มี query parameters ใดๆ (Criteria ว่าง)
	req := httptest.NewRequest(http.MethodGet, "/patient/search", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("คาดหวัง Status 400 แต่ได้ %d", rec.Code)
	}
}

func TestPatientHandler_SearchPatient_ResolverError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resolver := &mockResolver{Err: errors.New("unsupported hospital")}
	searchUC := application.NewSearchPatient(resolver)
	handler := httpDelivery.NewPatientHandler(searchUC)

	r := gin.New()
	r.GET("/patient/search", func(c *gin.Context) {
		c.Set("hospital", "UNKNOWN")
		handler.SearchPatient(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/patient/search?national_id=1100200300400", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("คาดหวัง Status 500 แต่ได้ %d", rec.Code)
	}
}

func TestSearchCriteria_Validate_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		criteria    domain.SearchCriteria
		expectError bool
	}{
		{
			name:        "Valid with NationalID only",
			criteria:    domain.SearchCriteria{NationalID: "1100200300400"},
			expectError: false,
		},
		{
			name:        "Valid with PassportID only",
			criteria:    domain.SearchCriteria{PassportID: "A12345678"},
			expectError: false,
		},
		{
			name:        "Valid with FirstName only",
			criteria:    domain.SearchCriteria{FirstName: "Somchai"},
			expectError: false,
		},
		{
			name:        "Valid with MiddleName only",
			criteria:    domain.SearchCriteria{MiddleName: "Ananda"},
			expectError: false,
		},
		{
			name:        "Valid with LastName only",
			criteria:    domain.SearchCriteria{LastName: "Jaidee"},
			expectError: false,
		},
		{
			name:        "Valid with DateOfBirth only",
			criteria:    domain.SearchCriteria{DateOfBirth: "1990-01-01"},
			expectError: false,
		},
		{
			name:        "Valid with PhoneNumber only",
			criteria:    domain.SearchCriteria{PhoneNumber: "0812345678"},
			expectError: false,
		},
		{
			name:        "Valid with Email only",
			criteria:    domain.SearchCriteria{Email: "somchai@example.com"},
			expectError: false,
		},
		{
			name: "Valid with Multiple Criteria",
			criteria: domain.SearchCriteria{
				FirstName: "Somchai",
				LastName:  "Jaidee",
				Email:     "somchai@example.com",
			},
			expectError: false,
		},
		{
			name:        "Invalid with All Empty Fields",
			criteria:    domain.SearchCriteria{}, // ว่างเปล่าทุก Field
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.criteria.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("[%s] คาดหวังว่าจะเกิด Error เนื่องจากไม่มีการระบุเงื่อนไขการค้นหา แต่ได้ nil", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("[%s] คาดหวังว่าผ่านการ Validate แต่เกิด Error: %v", tt.name, err)
				}
			}
		})
	}
}
