package infrastructure

import (
	"Hospital-Midderware/internal/domain"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTTokenProvider_GenerateToken_Success(t *testing.T) {
	secretKey := "my-secret-key-54321"
	provider := NewJWTTokenProvider(secretKey)

	staff, err := domain.CreateStaff("nattapong", "KhonNaRak5555", "HN11111")
	if err != nil {
		t.Fatalf("สร้าง staff ไม่สำเร็จ: %v", err)
	}

	tokenString, err := provider.GenerateToken(staff)

	if err != nil {
		t.Fatalf("คาดว่าสร้าง Token สำเร็จ แต่ได้ error: %v", err)
	}

	if tokenString == "" {
		t.Fatal("คาดว่าต้องได้ Token string แต่ได้ข้อความว่างเปล่าา")
	}
}

func TestJWTTokenProvider_GenerateToken_ShouldContainCorrectClaims(t *testing.T) {
	secretKey := "my-secret-key-9999"
	provider := NewJWTTokenProvider(secretKey)

	staff, _ := domain.CreateStaff("nattapong", "RakTukKhon8888", "HN6666")

	tokenString, err := provider.GenerateToken(staff)
	if err != nil {
		t.Fatalf("สร้าง Token ไม่สำเร็จ: %v", err)
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !parsedToken.Valid {
		t.Fatalf("Token ไม่ถูกต้อง หรือถูกแกะด้วย Secret Key ผิด: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("ไม่สามารถแปลง Claims เป็น jwt.MapClaims ได้")
	}

	if claims["username"] != "nattapong" {
		t.Errorf("คาดหวัง username %q แต่ได้ %v", "nattapong", claims["username"])
	}

	if claims["hospital_id"] != "HN6666" {
		t.Errorf("คาดหวัง hospital_id %q แต่ได้ %v", "HN6666", claims["hospital_id"])
	}

	if claims["exp"] == nil {
		t.Error("คาดว่าต้องมีวันหมดอายุ exp ใน token")
	}
}
