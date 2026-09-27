package domain_test

import (
	"Hospital-Middleware/internal/domain"
	"testing"
)

func TestVersion_Initial_Success(t *testing.T) {
	v := domain.InitialVersion()
	if v.CurrentNumber() != 1 {
		t.Errorf("คาดหวัง current_number เป็น 1 แต่ได้ %d", v.CurrentNumber())
	}
	if v.PreviousNumber() != 1 {
		t.Errorf("คาดหวัง previous_number เป็น 1 แต่ได้ %d", v.PreviousNumber())
	}
}

func TestVersion_Increment_Success(t *testing.T) {
	v := domain.InitialVersion()
	v2 := v.Increment()

	if v2.CurrentNumber() != 2 {
		t.Errorf("คาดหวัง current_number เป็น 2 แต่ได้ %d", v2.CurrentNumber())
	}
	if v2.PreviousNumber() != 1 {
		t.Errorf("คาดหวัง previous_number เป็น 1 แต่ได้ %d", v2.PreviousNumber())
	}

	v3 := v2.Increment()
	if v3.CurrentNumber() != 3 {
		t.Errorf("คาดหวัง current_number เป็น 3 แต่ได้ %d", v3.CurrentNumber())
	}
	if v3.PreviousNumber() != 2 {
		t.Errorf("คาดหวัง previous_number เป็น 2 แต่ได้ %d", v3.PreviousNumber())
	}
}

func TestVersion_NewVersion_ValidationErrors(t *testing.T) {
	// ทดสอบค่าน้อยกว่า 1
	_, err := domain.NewVersion(0, 1)
	if err == nil {
		t.Errorf("คาดหวังว่าจะได้ error เมื่อ current < 1 แต่กลับไม่เจอ error")
	}

	// ทดสอบ current < previous
	_, err = domain.NewVersion(1, 2)
	if err == nil {
		t.Errorf("คาดหวังว่าจะได้ error เมื่อ current < previous แต่กลับไม่เจอ error")
	}
}
