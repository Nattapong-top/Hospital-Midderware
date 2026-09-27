package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Hospital-Middleware/internal/domain"
	"Hospital-Middleware/internal/infrastructure"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func TestAuthMiddleware_Rejections(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const secret = "secret-key-6666"
	tokenProvider := infrastructure.NewJWTTokenProvider(secret)
	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, infrastructure.JWTClaims{
		Username: "Paa_Top_IT",
		Hospital: "HN99999",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("สร้าง expired token ไม่สำเร็จ: %v", err)
	}

	tests := []struct {
		name          string
		authorization string
	}{
		{name: "missing authorization header"},
		{name: "non-Bearer authorization scheme", authorization: "Basic abc"},
		{name: "fake token", authorization: "Bearer not-a-valid-token"},
		{name: "expired token", authorization: "Bearer " + expiredToken},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handlerCalled := false
			r := gin.New()
			r.Use(AuthMiddleware(tokenProvider))
			r.GET("/test-protected", func(c *gin.Context) {
				handlerCalled = true
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test-protected", nil)
			if test.authorization != "" {
				req.Header.Set("Authorization", test.authorization)
			}
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want %d; body: %s", rec.Code, http.StatusUnauthorized, rec.Body.String())
			}
			if handlerCalled {
				t.Error("protected handler was called for an invalid authorization")
			}
		})
	}
}
