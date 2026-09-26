package application

import (
	"Hospital-Midderware/internal/domain"
	"errors"
)

type AuthService struct {
	staffRepo     domain.StaffRepository
	hasher        domain.PasswordHasher
	tokenProvider domain.TokenProvider
}

func NewAuthService(
	staffRepo domain.StaffRepository,
	hasher domain.PasswordHasher,
	tokenProvider domain.TokenProvider,
) *AuthService {
	return &AuthService{
		staffRepo:     staffRepo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
	}

}

func (a *AuthService) Login(username, password string) (string, error) {
	staff, err := a.staffRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	isMatched := a.hasher.Compare(staff.Password.Value(), password)
	if !isMatched {
		return "", errors.New("invalid credentials")
	}

	token, err := a.tokenProvider.GenerateToken(staff)
	if err != nil {
		return "", err
	}

	return token, nil
}
