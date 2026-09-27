package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Hospital-Middleware/internal/application"
	"Hospital-Middleware/internal/domain"
	"Hospital-Middleware/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

// --- Mock Staff Repository ในชั้น Memory สำหรับ Test ---
type mockStaffRepo struct{}

func (m *mockStaffRepo) Create(ctx context.Context, staff *domain.Staff) error {
	//TODO implement me
	return nil
}

func (m *mockStaffRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	//TODO implement me
	return false, nil
}

func (m *mockStaffRepo) FindByUsername(username string) (*domain.Staff, error) {
	if username == "Paa_Top_IT" {
		staff, _ := domain.CreateStaff("Paa_Top_IT", "$2a$10$abcdefghijklmnopqrstuu", "HN99999")
		return staff, nil
	}
	return nil, nil
}

func (m *mockStaffRepo) Save(staff *domain.Staff) error {
	return nil
}

func TestMainRoutes_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Mock External Hospital A Server
	mockExternalHospA := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		patient := domain.PatientDTO{
			FirstNameTH: "สมชาย",
			NationalID:  "1100200300400",
			PatientHN:   "HN-12345",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(patient)
	}))
	defer mockExternalHospA.Close()

	// 2. Setup Mock Dependencies
	staffRepo := &mockStaffRepo{}
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider("test-secret-key")

	hospitalAAdapter := infrastructure.NewHospitalAAPIAdapter(mockExternalHospA.URL)
	hospitalResolver := infrastructure.NewHospitalResolver(hospitalAAdapter)

	authService := application.NewAuthService(staffRepo, hasher, tokenProvider)
	staffService := application.NewStaffService(staffRepo, hasher)
	searchPatientUseCase := application.NewSearchPatient(hospitalResolver)

	// 3. Setup Router
	router := setupRouter(authService, &staffService, searchPatientUseCase, tokenProvider)

	// --- Scenario 1: Access Protected Route Without Auth Header -> Should Fail 401 ---
	t.Run("GET /api/v1/patient/search without token -> 401 Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/patient/search?national_id=1100200300400", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("คาดหวัง Status %d แต่ได้ %d", http.StatusUnauthorized, rec.Code)
		}
	})

	// --- Scenario 2: Access Protected Route With Valid JWT Token -> Should Pass 200 ---
	t.Run("GET /api/v1/patient/search with valid token -> 200 OK", func(t *testing.T) {
		// สร้าง Token จำลอง
		staff, _ := domain.CreateStaff("Paa_Top_IT", "HashedPass123!", "HN99999")
		validToken, err := tokenProvider.GenerateToken(staff)
		if err != nil {
			t.Fatalf("สร้าง token ล้มเหลว: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/patient/search?national_id=1100200300400", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("คาดหวัง Status %d แต่ได้ %d (Response: %s)", http.StatusOK, rec.Code, rec.Body.String())
		}
	})
}
