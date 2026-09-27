package application_test

import (
	"context"
	"errors"
	"testing"

	"Hospital-Middleware/internal/application"
	"Hospital-Middleware/internal/domain"
)

// Mock Repository สำหรับ Staff Management
type mockStaffRepo struct {
	staffs              map[string]*domain.Staff
	existsError         error
	saveError           error
	findByUsernameError error
}

func newMockStaffRepo() *mockStaffRepo {
	return &mockStaffRepo{
		staffs: make(map[string]*domain.Staff),
	}
}

func (m *mockStaffRepo) FindByUsername(username string) (*domain.Staff, error) {
	if m.findByUsernameError != nil {
		return nil, m.findByUsernameError
	}
	staff, exists := m.staffs[username]
	if !exists {
		return nil, errors.New("staff not found")
	}
	return staff, nil
}

func (m *mockStaffRepo) Save(staff *domain.Staff) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.staffs[staff.Username.Value()] = staff
	return nil
}

func (m *mockStaffRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	if m.existsError != nil {
		return false, m.existsError
	}
	_, exists := m.staffs[username]
	return exists, nil
}

func (m *mockStaffRepo) Create(ctx context.Context, staff *domain.Staff) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.staffs[staff.Username.Value()] = staff
	return nil
}

// Mock PasswordHasher
type mockHasher struct {
	hashError error
}

func (m *mockHasher) Hash(password string) (string, error) {
	if m.hashError != nil {
		return "", m.hashError
	}
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

// Test 2: ทดสอบกรณี Username ซ้ำ
func TestStaffService_CreateStaff_DuplicateUsername(t *testing.T) {
	repo := newMockStaffRepo()
	hasher := &mockHasher{}
	service := application.NewStaffService(repo, hasher)

	ctx := context.Background()

	existingStaff, _ := domain.CreateStaff("Paa_Top_IT", "Password123!", "HN99999")
	_ = repo.Save(existingStaff)

	req := application.CreateStaffRequest{
		Username: "Paa_Top_IT",
		Password: "NewPassword123!",
		Hospital: "HN99999",
	}

	err := service.CreateStaff(ctx, req)

	if err == nil {
		t.Fatalf("คาดหวังว่าจะได้ Error เรื่อง Username ซ้ำ แต่กลับไม่เจอ Error")
	}
}

// Test 3: ทดสอบกรณีรหัสผ่านไม่ผ่านตามกฎ Domain (สั้นเกินไป)
func TestStaffService_CreateStaff_InvalidPassword(t *testing.T) {
	repo := newMockStaffRepo()
	hasher := &mockHasher{}
	service := application.NewStaffService(repo, hasher)

	ctx := context.Background()
	req := application.CreateStaffRequest{
		Username: "Paa_Top_IT",
		Password: "123", // สั้นเกินไป
		Hospital: "HN99999",
	}

	err := service.CreateStaff(ctx, req)
	if err == nil {
		t.Fatalf("คาดหวังว่าจะได้ Error เรื่องรหัสผ่าน แต่กลับไม่เจอ Error")
	}
}

// Test 4: ทดสอบกรณี Repository เกิดข้อผิดพลาดตอนเช็ค ExistsByUsername
func TestStaffService_CreateStaff_RepoExistsError(t *testing.T) {
	repo := newMockStaffRepo()
	repo.existsError = errors.New("database error")
	hasher := &mockHasher{}
	service := application.NewStaffService(repo, hasher)

	ctx := context.Background()
	req := application.CreateStaffRequest{
		Username: "Paa_Top_IT",
		Password: "Password123!",
		Hospital: "HN99999",
	}

	err := service.CreateStaff(ctx, req)
	if err == nil {
		t.Fatalf("คาดหวังว่าจะได้ Error จาก Repository แต่กลับไม่เจอ Error")
	}
}

// Test 5: ทดสอบกรณี Hasher เกิดข้อผิดพลาดตอน Hash รหัสผ่าน
func TestStaffService_CreateStaff_HashError(t *testing.T) {
	repo := newMockStaffRepo()
	hasher := &mockHasher{
		hashError: errors.New("hash failed"),
	}
	service := application.NewStaffService(repo, hasher)

	ctx := context.Background()
	req := application.CreateStaffRequest{
		Username: "Paa_Top_IT",
		Password: "Password123!",
		Hospital: "HN99999",
	}

	err := service.CreateStaff(ctx, req)
	if err == nil {
		t.Fatalf("คาดหวังว่าจะได้ Error จาก Hasher แต่กลับไม่เจอ Error")
	}
}

// Test 6: ทดสอบกรณี Repository เกิดข้อผิดพลาดตอน Save
func TestStaffService_CreateStaff_SaveError(t *testing.T) {
	repo := newMockStaffRepo()
	repo.saveError = errors.New("save failed")
	hasher := &mockHasher{}
	service := application.NewStaffService(repo, hasher)

	ctx := context.Background()
	req := application.CreateStaffRequest{
		Username: "Paa_Top_IT",
		Password: "Password123!",
		Hospital: "HN99999",
	}

	err := service.CreateStaff(ctx, req)
	if err == nil {
		t.Fatalf("คาดหวังว่าจะได้ Error ตอน Save แต่กลับไม่เจอ Error")
	}
}

// Test 7: ทดสอบกรณี Hospital ID ไม่ถูกต้องตาม Domain
func TestStaffService_CreateStaff_InvalidHospital(t *testing.T) {
	repo := newMockStaffRepo()
	hasher := &mockHasher{}
	service := application.NewStaffService(repo, hasher)

	ctx := context.Background()
	req := application.CreateStaffRequest{
		Username: "Paa_Top_IT",
		Password: "Password123!",
		Hospital: "", // Hospital ว่างหรือไม่ถูกต้อง
	}

	err := service.CreateStaff(ctx, req)
	if err == nil {
		t.Fatalf("คาดหวังว่าจะได้ Error เรื่อง Hospital ID แต่กลับไม่เจอ Error")
	}
}
