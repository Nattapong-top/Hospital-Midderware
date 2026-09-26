package application_test

import (
	"context"
	"errors"
	"testing"

	"Hospital-Midderware/internal/application"
	"Hospital-Midderware/internal/domain"
)

// Mock Repository สำหรับ Staff Management
type mockStaffRepo struct {
	staffs map[string]*domain.Staff
}

func newMockStaffRepo() *mockStaffRepo {
	return &mockStaffRepo{
		staffs: make(map[string]*domain.Staff),
	}
}

func (m *mockStaffRepo) FindByUsername(username string) (*domain.Staff, error) {
	staff, exists := m.staffs[username]
	if !exists {
		return nil, errors.New("staff not found")
	}
	return staff, nil
}

func (m *mockStaffRepo) Save(staff *domain.Staff) error {
	m.staffs[staff.Username.Value()] = staff
	return nil
}

func (m *mockStaffRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	_, exists := m.staffs[username]
	return exists, nil
}

func (m *mockStaffRepo) Create(ctx context.Context, staff *domain.Staff) error {
	m.staffs[staff.Username.Value()] = staff
	return nil
}

type mockHospitalRepo struct {
	validHospitals map[string]bool
}

func newMockHospitalRepo() *mockHospitalRepo {
	return &mockHospitalRepo{
		validHospitals: map[string]bool{
			"HN12345": true, // Hospital ที่มีจริงในระบบ
		},
	}
}

func (m *mockHospitalRepo) ExistsByID(ctx context.Context, id string) (bool, error) {
	return m.validHospitals[id], nil
}

// Mock PasswordHasher
type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) {
	return "hashed_" + password, nil
}

func (m *mockHasher) Compare(hashedPassword, password string) bool {
	return hashedPassword == "hashed_"+password
}

// Test 1: ทดสอบการสร้าง Staff สำเร็จ ข้อมูลต้องถูกเซฟลง Repository
func TestStaffService_CreateStaff_Success(t *testing.T) {
	repo := newMockStaffRepo()
	hasher := &mockHasher{}
	service := application.NewStaffService(repo, hasher)

	ctx := context.Background()
	req := application.CreateStaffRequest{
		Username: "Paa_Top_IT",
		Password: "Password123!",
		Hospital: "HN99999",
	}

	err := service.CreateStaff(ctx, req)
	if err != nil {
		t.Fatalf("คาดหวังว่าจะสร้าง Staff สำเร็จ แต่ได้ error: %v", err)
	}

	// ตรวจสอบว่าข้อมูลถูกบันทึกลง Repo จริง
	exists, _ := repo.ExistsByUsername(ctx, "Paa_Top_IT")
	if !exists {
		t.Errorf("คาดหวังว่าจะพบ Username Paa_Top_IT ใน Repository")
	}
}

func TestStaffService_CreateStaff_DuplicateUsername(t *testing.T) {
	repo := newMockStaffRepo()
	hasher := &mockHasher{}
	service := application.NewStaffService(repo, hasher)

	ctx := context.Background()

	existingStaff, _ := domain.CreateStaff("Paa_Top_IT", "hashed_Password123!", "HN99999")
	err := repo.Save(existingStaff)
	if err != nil {
		return
	}

	req := application.CreateStaffRequest{
		Username: "Paa_Top_IT",
		Password: "NewPassword123!",
		Hospital: "HN99999",
	}

	err = service.CreateStaff(ctx, req)

	// 3. ต้องได้ Error กลับมา
	if err == nil {
		t.Fatalf("คาดหวังว่าจะได้ Error เรื่อง Username ซ้ำ แต่กลับไม่เจอ Error")
	}
}
