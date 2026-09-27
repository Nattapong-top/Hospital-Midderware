package infrastructure_test

import (
	"testing"

	"Hospital-Middleware/internal/domain"
	"Hospital-Middleware/internal/infrastructure"
)

type dummyAdapter struct{}

func (d *dummyAdapter) Search(criteria domain.SearchCriteria) (*domain.PatientDTO, error) {
	return nil, nil
}

func TestHospitalResolver_Resolve_Success(t *testing.T) {
	adapter := &dummyAdapter{}
	resolver := infrastructure.NewHospitalResolver(adapter)

	supportedIDs := []string{"HOSP_A", "HN99999", "HN12345"}
	for _, id := range supportedIDs {
		t.Run("Resolve "+id, func(t *testing.T) {
			resolvedAdapter, err := resolver.Resolve(id)
			if err != nil {
				t.Errorf("คาดหวังว่าจะ resolve สำเร็จสำหรับ ID %s แต่ได้ error: %v", id, err)
			}
			if resolvedAdapter == nil {
				t.Errorf("คาดหวัง adapter แต่ได้ nil สำหรับ ID %s", id)
			}
		})
	}
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
