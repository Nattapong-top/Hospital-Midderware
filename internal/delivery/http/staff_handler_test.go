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

	"github.com/gin-gonic/gin"
)

// Helper ฟังก์ชันสำหรับสร้าง *gin.Context ในการทำ Unit Test
func setupGinTestContext(req *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	return c, rec
}

func TestStaffHandler_Login_Success(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider("secret-key-6666")

	hashedPassword, _ := hasher.Hash("KhonNaRak5555")
	staff, _ := domain.CreateStaff("Paa_Top_IT", hashedPassword, "HN99999")
	_ = repo.Save(staff)

	authService := application.NewAuthService(repo, hasher, tokenProvider)
	handler := NewStaffHandler(authService, nil)

	loginReq := LoginRequest{
		Username: "Paa_Top_IT",
		Password: "KhonNaRak5555",
		Hospital: "HN99999",
	}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(loginReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", &buf)
	req.Header.Set("Content-Type", "application/json")

	c, rec := setupGinTestContext(req)
	handler.Login(c)

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
	handler := NewStaffHandler(authService, nil)

	loginReq := LoginRequest{
		Username: "Paa_Top_IT",
		Password: "WrongPassword123",
		Hospital: "HN99999",
	}
	jsonBody, _ := json.Marshal(loginReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, rec := setupGinTestContext(req)
	handler.Login(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("คาดหวัง Status Code %d แต่ได้ %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestStaffHandler_Login_InvalidJSON_ShouldFail(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider("secret-key-6666")

	authService := application.NewAuthService(repo, hasher, tokenProvider)
	handler := NewStaffHandler(authService, nil)

	invalidJSON := []byte(`{"username": "Paa_Top_IT", "password":`)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	c, rec := setupGinTestContext(req)
	handler.Login(c)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("คาดหวัง Status Code %d แต่ได้ %d", http.StatusBadRequest, rec.Code)
	}
}

func TestStaffHandler_CreateStaff(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()

	staffService := application.NewStaffService(repo, hasher)
	handler := NewStaffHandler(nil, &staffService)

	createStaffReq := CreateStaffRequest{
		Username: "Paa_Top_IT",
		Password: "KhonNaRak5555",
		Hospital: "HN99999",
	}
	jsonBody, _ := json.Marshal(createStaffReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, rec := setupGinTestContext(req)
	handler.CreateStaff(c)

	if rec.Code != http.StatusCreated {
		t.Errorf("คาดหวัง Status Code %d แต่ได้ %d", http.StatusCreated, rec.Code)
	}
}

func TestStaffHandler_CreateStaff_InvalidJSON_ShouldFail(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()

	staffService := application.NewStaffService(repo, hasher)
	handler := NewStaffHandler(nil, &staffService)

	invalidJson := []byte(`{ username: "Paa_Top_IT", invalid_json... `)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff/create", bytes.NewBuffer(invalidJson))
	req.Header.Set("Content-Type", "application/json")

	c, rec := setupGinTestContext(req)
	handler.CreateStaff(c)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("คาดหวัง Status Code %d แต่ได้ %d", http.StatusBadRequest, rec.Code)
	}
}

func TestStaffHandler_CreateStaff_DuplicateUsername_ShouldFail(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()

	staffService := application.NewStaffService(repo, hasher)
	handler := NewStaffHandler(nil, &staffService)

	// 1. จำลองว่ามี Username "Paa_Top_IT" บันทึกอยู่ใน Repo แล้ว
	existingStaff, _ := domain.CreateStaff("Paa_Top_IT", "HashedPass123!", "HN99999")
	_ = repo.Save(existingStaff)

	// 2. พยายามส่ง Request สร้าง Username "Paa_Top_IT" ซ้ำ
	createStaffReq := CreateStaffRequest{
		Username: "Paa_Top_IT",
		Password: "NewPassword5555",
		Hospital: "HN99999",
	}
	jsonBody, _ := json.Marshal(createStaffReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, rec := setupGinTestContext(req)
	handler.CreateStaff(c)

	if rec.Code != http.StatusBadRequest && rec.Code != http.StatusConflict {
		t.Errorf("คาดหวัง Status Code 400 หรือ 409 แต่ได้ %d", rec.Code)
	}
}
