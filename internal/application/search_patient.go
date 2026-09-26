package application

import (
	"fmt"

	"Hospital-Midderware/internal/domain"
)

type SearchPatient struct {
	resolver domain.HospitalAPIResolver
}

func NewSearchPatient(resolver domain.HospitalAPIResolver) *SearchPatient {
	return &SearchPatient{
		resolver: resolver,
	}
}

func (uc *SearchPatient) Execute(hospitalID string, criteria domain.SearchCriteria) (*domain.SearchPatientResult, error) {
	// 1. Validate Search Criteria (ต้องมีอย่างน้อย 1 เงื่อนไขตาม Rule D05-14)
	if err := criteria.Validate(); err != nil {
		return nil, fmt.Errorf("invalid criteria: %w", err)
	}

	// 2. Resolve External Integration จาก hospital_id (ตาม Rule D05-11)
	adapter, err := uc.resolver.Resolve(hospitalID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve hospital API: %w", err)
	}

	// 3. ยิงค้นหาข้อมูลจาก External API
	patientDTO, err := adapter.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("external search failed: %w", err)
	}

	// 4. Return ผลลัพธ์กลับไป (ถ้าไม่เจอ patientDTO จะเป็น nil ตาม Rule D05-17)
	return &domain.SearchPatientResult{
		Patient: patientDTO,
	}, nil
}
