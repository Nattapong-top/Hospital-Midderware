package domain

import "testing"

func TestNewHospitalId_Success(t *testing.T) {
	hospitalId := "HN12345"
	hospitalIdVo, err := NewHospitalId(hospitalId)

	if err != nil {
		t.Fatalf("คาดว่าไม่มี error แต่มี error: %v", err)
	}

	if hospitalId != hospitalIdVo.Value() {
		t.Errorf("คาดหวังให้ hospitalId: %q, แต่ได้รับ %q",
			hospitalId,
			hospitalIdVo.Value(),
		)
	}
}

func TestNewHospitalId_entry(t *testing.T) {
	hospitalId := ""
	_, err := NewHospitalId(hospitalId)
	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อ hospitalID ว่าง")
	}
}

func TestNewHospitalId_ShouldTrimSpace(t *testing.T) {
	input := "  HN12345  "
	expected := "HN12345"

	hospitalIdVo, err := NewHospitalId(input)

	if err != nil {
		t.Fatalf("คาดว่าไม่มี error แต่มี error: %v", err)
	}

	if hospitalIdVo.Value() != expected {
		t.Errorf("คาดหวัง %q แต่ได้ %q", expected, hospitalIdVo.Value())
	}
}
