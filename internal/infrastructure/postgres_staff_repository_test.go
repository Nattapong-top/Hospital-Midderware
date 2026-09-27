package infrastructure

import (
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
		t.Fatalf("ping database ไม่ผ่าน (กรุณาเช็คว่า docker สตาร์ตแล้วหรือยัง): %v", err)
	}

	return db
}

func TestPostgresStaffRepository_SaveAndFindByUsername_Success(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		_, _ = db.Exec("DELETE FROM staffs WHERE username = $1", "Paa_Top_IT")
		_ = db.Close()
	})

	repo := NewPostgresStaffRepository(db)

	staff, err := domain.CreateStaff("Paa_Top_IT", "KhonNaRak5555", "HN99999")
	if err != nil {
		t.Fatalf("สร้าง staff ไม่สำเร็จ: %v", err)
	}

	err = repo.Save(staff)
	if err != nil {
		t.Fatalf("Save staff ลง Postgres ไม่สำเร็จ: %v", err)
	}

	foundStaff, err := repo.FindByUsername("Paa_Top_IT")
	if err != nil {
		t.Fatalf("FindByUsername ไม่สำเร็จ: %v", err)
	}

	if foundStaff.Username.Value() != "Paa_Top_IT" {
		t.Errorf("คาดหวัง username %q แต่ได้ %q", "Paa_Top_IT", foundStaff.Username.Value())
	}

	if foundStaff.HospitalId.Value() != "HN99999" {
		t.Errorf("คาดหวัง hospital_id %q แต่ได้ %q", "HN99999", foundStaff.HospitalId.Value())
	}
}

// Test Happy Path: ทดสอบการบันทึกและอัปเดต Version (Increment) ของ Staff ใน PostgreSQL
func TestPostgresStaffRepository_Version_HappyPath(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() {
		_, _ = db.Exec("DELETE FROM staffs WHERE username = $1", "Version_Staff_Test")
		_ = db.Close()
	})

	repo := NewPostgresStaffRepository(db)

	// 1. สร้าง Staff ใหม่ (Version เริ่มต้น 1, 1)
	staff, err := domain.CreateStaff("Version_Staff_Test", "Password123!", "HN99999")
	if err != nil {
		t.Fatalf("สร้าง staff ไม่สำเร็จ: %v", err)
	}

	if staff.Version.CurrentNumber() != 1 || staff.Version.PreviousNumber() != 1 {
		t.Fatalf("เวอร์ชันเริ่มต้นไม่ถูกต้อง")
	}

	err = repo.Save(staff)
	if err != nil {
		t.Fatalf("Save staff ครั้งแรกไม่สำเร็จ: %v", err)
	}

	// 2. โหลดขึ้นมาเช็คว่า Version เป็น 1, 1 จริง
	loaded1, err := repo.FindByUsername("Version_Staff_Test")
	if err != nil {
		t.Fatalf("FindByUsername ไม่สำเร็จ: %v", err)
	}
	if loaded1.Version.CurrentNumber() != 1 || loaded1.Version.PreviousNumber() != 1 {
		t.Errorf("คาดหวัง version 1/1 แต่ได้ %d/%d", loaded1.Version.CurrentNumber(), loaded1.Version.PreviousNumber())
	}

	// 3. ทำการ Increment Version (เป็น 2, 1) และ Save ทับ
	newVersion := loaded1.Version.Increment()
	updatedStaff, err := domain.RestoreStaff("Version_Staff_Test", "Password123!", "HN99999", newVersion.CurrentNumber(), newVersion.PreviousNumber())
	if err != nil {
		t.Fatalf("RestoreStaff ไม่สำเร็จ: %v", err)
	}

	err = repo.Save(updatedStaff)
	if err != nil {
		t.Fatalf("Save staff ครั้งที่สอง (หลัง increment) ไม่สำเร็จ: %v", err)
	}

	// 4. โหลดขึ้นมาเช็คว่า Version อัปเดตเป็น 2, 1 จริง
	loaded2, err := repo.FindByUsername("Version_Staff_Test")
	if err != nil {
		t.Fatalf("FindByUsername ไม่สำเร็จ: %v", err)
	}
	if loaded2.Version.CurrentNumber() != 2 || loaded2.Version.PreviousNumber() != 1 {
		t.Errorf("คาดหวัง version 2/1 แต่ได้ %d/%d", loaded2.Version.CurrentNumber(), loaded2.Version.PreviousNumber())
	}
}

// Test Sad Path: ทดสอบกรณีค้นหา Username ที่ไม่มีอยู่จริง และกรณี Version ไม่ถูกต้อง
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
