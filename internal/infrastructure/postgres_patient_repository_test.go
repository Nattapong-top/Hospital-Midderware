package infrastructure

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgresPatientsTable_InsertAndQuery_Success(t *testing.T) {
	connStr := "host=127.0.0.1 port=5432 user=postgres password=KhonNaRak5555 dbname=hospital_middleware sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skipf("ข้าม Integration Test เนื่องจากไม่สามารถเชื่อมต่อ Database ได้: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("ข้าม Integration Test เนื่องจาก Database Ping ไม่ผ่าน: %v", err)
	}

	// 1. Insert Test Patient
	insertQuery := `
		INSERT INTO patients (patient_hn, hospital_id, national_id, first_name_th, last_name_th)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (patient_hn) DO NOTHING
	`
	_, err = db.Exec(insertQuery, "HN-TEST-999", "HN99999", "1100200300999", "ทดสอบ", "ระบบ")
	if err != nil {
		t.Fatalf("insert test patient ไม่สำเร็จ: %v", err)
	}

	// 2. Query Test Patient
	selectQuery := `SELECT patient_hn, national_id, first_name_th, last_name_th FROM patients WHERE patient_hn = $1`
	var hn, nationalID, firstName, lastName string
	err = db.QueryRow(selectQuery, "HN-TEST-999").Scan(&hn, &nationalID, &firstName, &lastName)
	if err != nil {
		t.Fatalf("query test patient ไม่สำเร็จ: %v", err)
	}

	if hn != "HN-TEST-999" || nationalID != "1100200300999" || firstName != "ทดสอบ" {
		t.Errorf("ข้อมูลผู้ป่วยที่ดึงได้ไม่ตรงกับที่ insert: got %s, %s, %s", hn, nationalID, firstName)
	}

	// 3. Cleanup
	_, _ = db.Exec("DELETE FROM patients WHERE patient_hn = $1", "HN-TEST-999")
}
