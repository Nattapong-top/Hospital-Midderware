package domain

import "testing"

func TestCreateStaff_Success(t *testing.T) {
	username := "nattapong"
	password := "khonnarak"
	hospitalId := "HN99999"

	staff, err := CreateStaff(username, password, hospitalId)

	if err != nil {
		t.Fatalf("คาดว่าไม่มี error แต่พบ error: %v", err)
	}

	if staff == nil {
		t.Fatal("คาดว่าจะสร้าง staff สำเร็จ แต่ได้ค่า nil")
	}

	if staff.Username.Value() != username {
		t.Errorf("คาดหวัง username %q, got %q", username, staff.Username.Value())
	}

	if staff.Password.Value() != password {
		t.Errorf("คาดหวัง password %q, got %q", password, staff.Password.Value())
	}

	if staff.HospitalId.Value() != hospitalId {
		t.Errorf("คาดหวัง HospitalId %q, got %q", hospitalId, staff.HospitalId.Value())
	}
}

func TestCreateStaff_InvalidUsername_shouldFail(t *testing.T) {
	username := "      "
	password := "khonnarak"
	hospitalId := "HN99999"

	_, err := CreateStaff(username, password, hospitalId)
	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อ username ว่าง")

	}
}

func TestCreateStaff_InvalidPassword_shouldFail(t *testing.T) {
	username := "nattapong"
	password := "     "
	hospitalId := "HN99999"
	_, err := CreateStaff(username, password, hospitalId)
	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อ password ว่าง")
	}
}

func TestCreateStaff_InvalidHospitalId_shouldFail(t *testing.T) {
	username := "nattapong"
	password := "khonnarak"
	hospitalId := ""

	_, err := CreateStaff(username, password, hospitalId)
	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อ hospitalId ว่าง")
	}
}
