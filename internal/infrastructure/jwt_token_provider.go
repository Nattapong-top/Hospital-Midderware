package infrastructure

import (
	"Hospital-Midderware/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenProvider struct {
	secretKey []byte
}

func NewJWTTokenProvider(secretKey string) *JwtTokenProvider {
	return &JwtTokenProvider{
		secretKey: []byte(secretKey),
	}
}

func (j *JwtTokenProvider) GenerateToken(staff *domain.Staff) (string, error) {

	if staff == nil {
		return "", errors.New("ไม่มีข้อมูล staff ในการสร้างบัตรผ่านครับ")
	}

	claims := jwt.MapClaims{
		"username":    staff.Username.Value(),
		"hospital_id": staff.HospitalId.Value(),
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}
