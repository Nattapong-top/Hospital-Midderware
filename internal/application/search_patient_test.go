package application_test

import (
	"testing"

	"Hospital-Midderware/internal/application"
	"Hospital-Midderware/internal/domain"
)

// Mock Adapter สำหรับ Unit Test
type mockAdapter struct {
	patient *domain.PatientDTO
	err     error
}

func (m *mockAdapter) Search(criteria domain.SearchCriteria) (*domain.PatientDTO, error) {
	return m.patient, m.err
}

// Mock Resolver สำหรับ Unit Test
type mockResolver struct {
	adapter domain.ExternalAPIAdapter
	err     error
}

func (m *mockResolver) Resolve(hospitalID string) (domain.ExternalAPIAdapter, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.adapter, nil
}

func TestSearchPatient_Execute_Success(t *testing.T) {
	expectedPatient := &domain.PatientDTO{
		FirstNameTH: "สมชาย",
		NationalID:  "1100200300400",
		PatientHN:   "HN-12345",
	}

	resolver := &mockResolver{
		adapter: &mockAdapter{patient: expectedPatient},
	}

	useCase := application.NewSearchPatient(resolver)
	criteria := domain.SearchCriteria{NationalID: "1100200300400"}

	result, err := useCase.Execute("HOSP_A", criteria)

	if err != nil {
		t.Fatalf("คาดหวังว่าค้นหาสำเร็จ แต่ได้ error: %v", err)
	}
	if result.Patient == nil {
		t.Fatal("คาดหวังผลลัพธ์ข้อมูลผู้ป่วย แต่ได้ nil")
	}
	if result.Patient.FirstNameTH != "สมชาย" {
		t.Errorf("คาดหวังชื่อ 'สมชาย' แต่ได้ %s", result.Patient.FirstNameTH)
	}
}

func TestSearchPatient_Execute_InvalidCriteria_ShouldFail(t *testing.T) {
	resolver := &mockResolver{}
	useCase := application.NewSearchPatient(resolver)

	// ส่ง criteria ว่างเปล่า
	_, err := useCase.Execute("HOSP_A", domain.SearchCriteria{})

	if err == nil {
		t.Error("คาดหวัง error เนื่องจากไม่มีการระบุเงื่อนไขการค้นหา")
	}
}
