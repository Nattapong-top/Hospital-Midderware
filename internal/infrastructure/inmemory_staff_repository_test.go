package infrastructure

import (
	"Hospital-Midderware/internal/domain"
	"testing"
)

func TestInMemoryStaffRepository_FindByUsername_NotFound_ShouldFail(t *testing.T) {
	repo := NewInMemoryStaffRepository()

	_, err := repo.FindByUsername("unknown_username")

	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อหาพนักงานไม่เจอ")
	}
}

func TestInMemoryStaffRepository_SaveAndFindByUsername_Success(t *testing.T) {
	repo := NewInMemoryStaffRepository()

	staff, err := domain.CreateStaff("nattapong", "khonnarak5555", "HN99999")
	if err != nil {
		t.Fatalf("สร้าง staff ไม่สำเร็จ: %v", err)
	}

	err = repo.Save(staff)
	if err != nil {
		t.Fatalf("บันทึก staff ไม่สำเร็จ: %v", err)
	}

	foundStaff, err := repo.FindByUsername("nattapong")
	if err != nil {
		t.Fatalf("ค้นหา staff ไม่พบ: %v", err)
	}

	if foundStaff.Username.Value() != staff.Username.Value() {
		t.Errorf("คาดหวัง username %q แต่ได้ %qq", staff.Username.Value(), foundStaff.Username.Value())
	}
}
