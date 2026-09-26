package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Hospital-Midderware/internal/application"
	"Hospital-Midderware/internal/domain"
	"Hospital-Midderware/internal/infrastructure"
)

func TestStaffHandler_Login_Success(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider("secret-key-6666")

	hashedPassword, _ := hasher.Hash("KhonNaRak5555")
	staff, _ := domain.CreateStaff("Paa_Top_IT", hashedPassword, "HN99999")
	_ = repo.Save(staff)

	authService := application.NewAuthService(repo, hasher, tokenProvider)
	handler := NewStaffHandler(authService)

	loginReq := LoginRequest{
		Username: "Paa_Top_IT",
		Password: "KhonNaRak5555",
		Hospital: "HN99999",
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(loginReq); err != nil {
		t.Fatalf("ไม่สามารถ encode request body ได้: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", &buf)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("คาดหวัง Status Code %d แต่ได้ %d", http.StatusOK, rec.Code)
	}

	var res LoginResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("ไม่สามารถ parse response body ได้: %v", err)
	}

	if res.Token == "" {
		t.Error("คาดว่าต้องได้รับ Token แต่ได้ข้อความว่างเปล่า")
	}
}

func TestStaffHandler_Login_InvalidCredentials_ShouldFail(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider("secret-key-6666")

	hashedPassword, _ := hasher.Hash("KhonNaRak5555")
	staff, _ := domain.CreateStaff("Paa_Top_IT", hashedPassword, "HN99999")
	_ = repo.Save(staff)

	authService := application.NewAuthService(repo, hasher, tokenProvider)
	handler := NewStaffHandler(authService)

	loginReq := LoginRequest{
		Username: "Paa_Top_IT",
		Password: "WrongPassword123",
		Hospital: "HN99999",
	}
	jsonBody, _ := json.Marshal(loginReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("คาดหวัง Status Code %d แต่ได้ %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestStaffHandler_Login_InvalidJSON_ShouldFail(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider("secret-key-6666")

	authService := application.NewAuthService(repo, hasher, tokenProvider)
	handler := NewStaffHandler(authService)

	invalidJSON := []byte(`{"username": "Paa_Top_IT", "password":`)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("คาดหวัง Status Code %d แต่ได้ %d", http.StatusBadRequest, rec.Code)
	}
}
