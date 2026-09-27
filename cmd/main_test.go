package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"Hospital-Middleware/internal/application"
	"Hospital-Middleware/internal/domain"
	"Hospital-Middleware/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

// --- Mock Staff Repository ในชั้น Memory สำหรับ Test ---
type mockStaffRepo struct {
	staffs map[string]*domain.Staff
}

func (m *mockStaffRepo) Create(ctx context.Context, staff *domain.Staff) error {
	if m.staffs == nil {
		m.staffs = make(map[string]*domain.Staff)
	}
	if _, exists := m.staffs[staff.Username.Value()]; exists {
		return errors.New("username นี้มีอยู่ในระบบแล้วครับ")
	}
	m.staffs[staff.Username.Value()] = staff
	return nil
}

func (m *mockStaffRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	if m.staffs == nil {
		m.staffs = make(map[string]*domain.Staff)
	}
	_, exists := m.staffs[username]
	return exists, nil
}

func (m *mockStaffRepo) FindByUsername(username string) (*domain.Staff, error) {
	if username == "Paa_Top_IT" {
		hashed, _ := infrastructure.NewBcryptHasher().Hash("Password123!")
		staff, _ := domain.CreateStaff("Paa_Top_IT", hashed, "HN99999")
		return staff, nil
	}
	if m.staffs != nil {
		if staff, exists := m.staffs[username]; exists {
			return staff, nil
		}
	}
	return nil, errors.New("ไม่พบข้อมูลพนักงานในระบบครับ")
}

func (m *mockStaffRepo) Save(staff *domain.Staff) error {
	if m.staffs == nil {
		m.staffs = make(map[string]*domain.Staff)
	}
	m.staffs[staff.Username.Value()] = staff
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
	staffRepo := &mockStaffRepo{staffs: make(map[string]*domain.Staff)}
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider("test-secret-key")

	hospitalAAdapter := infrastructure.NewHospitalAAPIAdapter(mockExternalHospA.URL)
	hospitalResolver := infrastructure.NewHospitalResolver(hospitalAAdapter)

	authService := application.NewAuthService(staffRepo, hasher, tokenProvider)
	staffService := application.NewStaffService(staffRepo, hasher)
	searchPatientUseCase := application.NewSearchPatient(hospitalResolver)

	// 3. Setup Router
	router := setupRouter(authService, &staffService, searchPatientUseCase, tokenProvider)

	// --- Test /staff/create (Success) ---
	t.Run("POST /staff/create Success", func(t *testing.T) {
		body := map[string]string{
			"username": "new_staff",
			"password": "Password123!",
			"hospital": "HN99999",
		}
		jsonBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/staff/create", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("คาดหวัง Status 201 แต่ได้ %d (Response: %s)", rec.Code, rec.Body.String())
		}
	})

	// --- Test /staff/login (Success) ---
	t.Run("POST /staff/login Success", func(t *testing.T) {
		// Create staff first
		hashed, _ := hasher.Hash("Password123!")
		staff, _ := domain.CreateStaff("login_staff", hashed, "HN99999")
		_ = staffRepo.Save(staff)

		body := map[string]string{
			"username": "login_staff",
			"password": "Password123!",
			"hospital": "HN99999",
		}
		jsonBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("คาดหวัง Status 200 แต่ได้ %d (Response: %s)", rec.Code, rec.Body.String())
		}
	})

	// --- Test /staff/login (Unauthorized / Wrong Password) ---
	t.Run("POST /staff/login Wrong Password -> 401", func(t *testing.T) {
		body := map[string]string{
			"username": "login_staff",
			"password": "WrongPassword!",
			"hospital": "HN99999",
		}
		jsonBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("คาดหวัง Status 401 แต่ได้ %d (Response: %s)", rec.Code, rec.Body.String())
		}
	})

	// --- Scenario 1: Access Protected Route Without Auth Header -> Should Fail 401 ---
	t.Run("GET /patient/search without token -> 401 Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/patient/search?national_id=1100200300400", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("คาดหวัง Status %d แต่ได้ %d", http.StatusUnauthorized, rec.Code)
		}
	})

	// --- Scenario 2: Access Protected Route With Valid JWT Token -> Should Pass 200 ---
	t.Run("GET /patient/search with valid token -> 200 OK", func(t *testing.T) {
		staff, _ := domain.CreateStaff("Paa_Top_IT", "HashedPass123!", "HN99999")
		validToken, err := tokenProvider.GenerateToken(staff)
		if err != nil {
			t.Fatalf("สร้าง token ล้มเหลว: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/patient/search?national_id=1100200300400", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("คาดหวัง Status %d แต่ได้ %d (Response: %s)", http.StatusOK, rec.Code, rec.Body.String())
		}
	})
}
