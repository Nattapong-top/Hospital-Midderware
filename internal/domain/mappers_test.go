package domain_test

import (
	"testing"

	"Hospital-Middleware/internal/domain"
)

func TestPatientToDTO(t *testing.T) {
	dto := domain.PatientToDTO("HN001", "1100200300400", "", "สมชาย", "ใจดี", "Somchai", "Jaidee", "M", "1990-01-01", "0812345678", "somchai@example.com")
	if dto.PatientHN != "HN001" || dto.NationalID != "1100200300400" || dto.FirstNameTH != "สมชาย" {
		t.Errorf("PatientToDTO mapping failed: got %+v", dto)
	}
}

func TestStaffToSummary(t *testing.T) {
	hashedPassword := "$2a$10$dummyhashedpassword"
	staff, err := domain.CreateStaff("test_staff", hashedPassword, domain.HospitalIDBangkok)
	if err != nil {
		t.Fatalf("unexpected error creating staff: %v", err)
	}

	summary := domain.StaffToSummary(staff)
	if summary.Username != "test_staff" || summary.HospitalID != domain.HospitalIDBangkok {
		t.Errorf("StaffToSummary mapping failed: got %+v", summary)
	}

	emptySummary := domain.StaffToSummary(nil)
	if emptySummary.Username != "" {
		t.Error("expected empty summary for nil staff")
	}
}
