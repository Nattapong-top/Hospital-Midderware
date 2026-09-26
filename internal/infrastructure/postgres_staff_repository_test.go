package infrastructure

import (
	"Hospital-Midderware/internal/domain"
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

	_, _ = db.Exec("DELETE FROM staffs WHERE username = $1", "Paa_Top_IT")
}
