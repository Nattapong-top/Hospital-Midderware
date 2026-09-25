package infrastructure

import (
	"Hospital-Midderware/internal/domain"
	"errors"
)

type InMemoryStaffRepository struct {
	staffs map[string]*domain.Staff
}

func NewInMemoryStaffRepository() *InMemoryStaffRepository {
	return &InMemoryStaffRepository{
		staffs: make(map[string]*domain.Staff),
	}
}

func (r *InMemoryStaffRepository) FindByUsername(username string) (*domain.Staff, error) {
	staff, exists := r.staffs[username]
	if !exists {
		return nil, errors.New("ไม่พบข้อมูลพนักงานในระบบ")
	}
	return staff, nil
}

func (r *InMemoryStaffRepository) Save(staff *domain.Staff) error {
	r.staffs[staff.Username.Value()] = staff
	return nil
}
