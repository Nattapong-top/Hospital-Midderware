package infrastructure

import (
	"errors"
	"time"

	"Hospital-Midderware/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Username string `json:"username"`
	Hospital string `json:"hospital_id"`
	jwt.RegisteredClaims
}

type JWTTokenProvider struct {
	secretKey string
}

func NewJWTTokenProvider(secretKey string) *JWTTokenProvider {
	return &JWTTokenProvider{secretKey: secretKey}
}

func (p *JWTTokenProvider) GenerateToken(staff *domain.Staff) (string, error) {

	if staff == nil {
		return "", errors.New("cannot generate token for nil staff")
	}

	claims := JWTClaims{
		Username:  staff.Username.Value(),
		Hospital:  staff.HospitalId.Value(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.secretKey))
}

func (p *JWTTokenProvider) ValidateToken(tokenStr string) (*domain.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(p.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &domain.CustomClaims{
		Username: claims.Username,
		Hospital: claims.Hospital,
	}, nil
}
