package application

import (
	"Hospital-Midderware/internal/domain"
	"context"
	"errors"
)

type CreateStaffRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

type StaffService struct {
	staffRepo domain.StaffRepository
	hasher    domain.PasswordHasher
}

func NewStaffService(staffRepo domain.StaffRepository, hasher domain.PasswordHasher) StaffService {
	return StaffService{
		staffRepo: staffRepo,
		hasher:    hasher,
	}
}

func (s *StaffService) CreateStaff(ctx context.Context, req CreateStaffRequest) error {

	exists, err := s.staffRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username นีัมีอยู่ในระบบแล้วครับ")
	}

	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		return err
	}

	staff, err := domain.CreateStaff(req.Username, hashedPassword, req.Hospital)
	if err != nil {
		return err
	}

	return s.staffRepo.Save(staff)
}
