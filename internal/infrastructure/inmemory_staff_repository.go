package infrastructure

import (
	"Hospital-Middleware/internal/domain"
	"context"
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
		return nil, errors.New("ไม่พบข้อมูลพนักงานในระบบครับ")
	}
	return staff, nil
}

func (r *InMemoryStaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	username := staff.Username.Value()
	if _, exists := r.staffs[username]; exists {
		return domain.ErrUserAlreadyExists
	}
	r.staffs[username] = staff
	return nil
}

func (r *InMemoryStaffRepository) Save(staff *domain.Staff) error {
	username := staff.Username.Value()
	existing, exists := r.staffs[username]
	if !exists {
		r.staffs[username] = staff
		return nil
	}

	// Optimistic Locking check
	if existing.Version.CurrentNumber() != staff.Version.PreviousNumber() {
		return errors.New("ข้อมูลพนักงานถูกแก้ไขโดยผู้อื่นแล้ว กรุณาลองใหม่อีกครั้งครับ (Optimistic Lock Conflict)")
	}

	r.staffs[username] = staff
	return nil
}

func (r *InMemoryStaffRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	_, exists := r.staffs[username]
	return exists, nil
}
