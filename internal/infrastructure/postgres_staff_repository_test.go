package infrastructure

import (
	"context"
	"Hospital-Middleware/internal/domain"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	connStr := "host=127.0.0.1 port=5432 user=postgres password=KhonNaRak5555 dbname=hospital_middleware sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("ไม่สามารถเชื่อมต่อ database ได้: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Skipf("ข้าม Integration Test เนื่องจากไม่สามารถเชื่อมต่อ Database ได้ (กรุณาเช็ค Docker): %v", err)
	}

	return db
}

func TestPostgresStaffRepository_CreateAndSave_Success(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	t.Cleanup(func() {
		_, _ = db.Exec("DELETE FROM staffs WHERE username = $1", "OpLock_Staff")
		_ = db.Close()
	})

	repo := NewPostgresStaffRepository(db)

	// 1. Test Create (Insert)
	staff, err := domain.CreateStaff("OpLock_Staff", "Password123!", "HN99999")
	if err != nil {
		t.Fatalf("สร้าง staff ไม่สำเร็จ: %v", err)
	}

	err = repo.Create(ctx, staff)
	if err != nil {
		t.Fatalf("Create staff ลง Postgres ไม่สำเร็จ: %v", err)
	}

	// 2. Test FindByUsername
	found, err := repo.FindByUsername("OpLock_Staff")
	if err != nil {
		t.Fatalf("FindByUsername ไม่สำเร็จ: %v", err)
	}
	if found.Version.CurrentNumber() != 1 {
		t.Errorf("คาดหวัง version 1 แต่ได้ %d", found.Version.CurrentNumber())
	}

	// 3. Test Save with Optimistic Locking (Update version 1 -> 2)
	newVersion := found.Version.Increment()
	updatedStaff, err := domain.RestoreStaff("OpLock_Staff", "NewPassword123!", "HN99999", newVersion.CurrentNumber(), newVersion.PreviousNumber())
	if err != nil {
		t.Fatalf("RestoreStaff ไม่สำเร็จ: %v", err)
	}

	err = repo.Save(updatedStaff)
	if err != nil {
		t.Fatalf("Save (Optimistic Update) ไม่สำเร็จ: %v", err)
	}

	// 4. Test Optimistic Lock Conflict (พยายาม Update ด้วยเวอร์ชันเก่าที่หมดอายุแล้ว)
	staleStaff, err := domain.RestoreStaff("OpLock_Staff", "AnotherPassword!", "HN99999", 2, 1) // version เก่า (current=2, previous=1 ซึ่งปัจจุบัน DB เป็น current=2 ไปแล้ว)
	if err != nil {
		t.Fatalf("RestoreStaff stale ไม่สำเร็จ: %v", err)
	}

	err = repo.Save(staleStaff)
	if err == nil {
		t.Fatalf("คาดหวังว่าจะเกิด Optimistic Lock Conflict error แต่ไม่มี error เกิดขึ้น")
	}
}

func TestPostgresStaffRepository_Version_SadPath(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		_ = db.Close()
	})

	repo := NewPostgresStaffRepository(db)

	// 1. ค้นหา username ที่ไม่มีในระบบ ต้องได้ Error ภาษาไทย
	_, err := repo.FindByUsername("NonExistentUser12345")
	if err == nil {
		t.Fatalf("คาดหวังว่าจะได้ Error เมื่อค้นหา username ที่ไม่มี แต่ไม่เจอ Error")
	}

	expectedErrMsg := "ไม่พบข้อมูลพนักงานในระบบครับ"
	if err.Error() != expectedErrMsg {
		t.Errorf("คาดหวัง error message %q แต่ได้ %q", expectedErrMsg, err.Error())
	}

	// 2. ทดสอบ RestoreStaff ด้วย Version ที่ไม่ถูกต้อง (เช่น current < previous)
	_, err = domain.RestoreStaff("TestUser", "Password123!", "HN99999", 1, 2)
	if err == nil {
		t.Fatalf("คาดหวังว่าจะได้ Error เมื่อ version ไม่ถูกต้อง แต่ไม่เจอ Error")
	}
}
