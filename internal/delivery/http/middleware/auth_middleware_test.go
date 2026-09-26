package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"Hospital-Midderware/internal/domain"
	"Hospital-Midderware/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tokenProvider := infrastructure.NewJWTTokenProvider("secret-key-6666")

	// 1. สร้าง Staff Entity สำหรับเจน Token
	staff, _ := domain.CreateStaff("Paa_Top_IT", "HashedPass123!", "HN99999")

	// 2. เจน Token โดยส่ง *domain.Staff
	validToken, err := tokenProvider.GenerateToken(staff)
	if err != nil {
		t.Fatalf("ไม่สามารถสร้าง Token สำหรับทดสอบได้: %v", err)
	}

	// 3. ตั้งค่า Router และ Middleware
	r := gin.New()
	r.Use(AuthMiddleware(tokenProvider))

	var capturedUsername, capturedHospital string
	r.GET("/test-protected", func(c *gin.Context) {
		capturedUsername = c.GetString("username")
		capturedHospital = c.GetString("hospital")
		c.Status(http.StatusOK)
	})

	// 4. ยิง Request พร้อม Header
	req := httptest.NewRequest(http.MethodGet, "/test-protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	// 5. ตรวจสอบผลลัพธ์
	if rec.Code != http.StatusOK {
		t.Errorf("คาดหวัง Status Code %d แต่ได้ %d", http.StatusOK, rec.Code)
	}
	if capturedUsername != "Paa_Top_IT" {
		t.Errorf("คาดหวัง username %s แต่ได้ %s", "Paa_Top_IT", capturedUsername)
	}
	if capturedHospital != "HN99999" {
		t.Errorf("คาดหวัง hospital %s แต่ได้ %s", "HN99999", capturedHospital)
	}
}
