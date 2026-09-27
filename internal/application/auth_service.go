package application

import (
	"Hospital-Middleware/internal/domain"
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

func (a *AuthService) Login(username, password, hospitalId string) (string, error) {
	staff, err := a.staffRepo.FindByUsername(username)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if staff.HospitalId.Value() != hospitalId {
		return "", domain.ErrInvalidCredentials
	}

	isMatched := a.hasher.Compare(staff.Password.Value(), password)
	if !isMatched {
		return "", domain.ErrInvalidCredentials
	}

	token, err := a.tokenProvider.GenerateToken(staff)
	if err != nil {
		return "", err
	}

	return token, nil
}
