package infrastructure_test

import (
	"testing"

	"Hospital-Middleware/internal/domain"
	"Hospital-Middleware/internal/infrastructure"
)

type dummyAdapter struct {
	name string
}

func (d *dummyAdapter) Search(criteria domain.SearchCriteria) (*domain.PatientDTO, error) {
	return nil, nil
}

func TestHospitalResolver_Resolve_Success(t *testing.T) {
	adapterA := &dummyAdapter{name: "AdapterA"}
	adapterB := &dummyAdapter{name: "AdapterB"}
	resolver := infrastructure.NewHospitalResolver(adapterA, adapterB)

	t.Run("Resolve HN12345 to Hospital A", func(t *testing.T) {
		resolved, err := resolver.Resolve("HN12345")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resolved != adapterA {
			t.Errorf("expected adapterA for HN12345")
		}
	})

	t.Run("Resolve HOSP_A to Hospital A", func(t *testing.T) {
		resolved, err := resolver.Resolve("HOSP_A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resolved != adapterA {
			t.Errorf("expected adapterA for HOSP_A")
		}
	})

	t.Run("Resolve HN99999 to Hospital B", func(t *testing.T) {
		resolved, err := resolver.Resolve("HN99999")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resolved != adapterB {
			t.Errorf("expected adapterB for HN99999")
		}
	})

	t.Run("Resolve HOSP_B to Hospital B", func(t *testing.T) {
		resolved, err := resolver.Resolve("HOSP_B")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resolved != adapterB {
			t.Errorf("expected adapterB for HOSP_B")
		}
	})
}

func TestHospitalResolver_Resolve_UnknownID_ShouldFailWithThaiError(t *testing.T) {
	adapter := &dummyAdapter{}
	resolver := infrastructure.NewHospitalResolver(adapter)

	_, err := resolver.Resolve("UNKNOWN_HOSPITAL")
	if err == nil {
		t.Fatalf("คาดหวัง error เมื่อ resolve รหัสโรงพยาบาลที่ไม่รู้จัก แต่กลับไม่เจอ error")
	}

	expectedErrorMsg := "ไม่พบรหัสโรงพยาบาลนี้ในระบบ หรือยังไม่รองรับครับ"
	if err.Error() != expectedErrorMsg {
		t.Errorf("คาดหวัง error message %q แต่ได้ %q", expectedErrorMsg, err.Error())
	}
}

func TestHospitalResolver_Resolve_EmptyID_ShouldFail(t *testing.T) {
	adapter := &dummyAdapter{}
	resolver := infrastructure.NewHospitalResolver(adapter)

	_, err := resolver.Resolve("")
	if err == nil {
		t.Fatalf("คาดหวัง error เมื่อ resolve รหัสโรงพยาบาลค่าว่าง แต่กลับไม่เจอ error")
	}

	expectedErrorMsg := "ไม่พบรหัสโรงพยาบาลนี้ในระบบ หรือยังไม่รองรับครับ"
	if err.Error() != expectedErrorMsg {
		t.Errorf("คาดหวัง error message %q แต่ได้ %q", expectedErrorMsg, err.Error())
	}
}
