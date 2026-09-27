package infrastructure

import (
	"errors"

	"Hospital-Middleware/internal/domain"
)

type HospitalResolver struct {
	hospitalAAdapter domain.ExternalAPIAdapter
}

func NewHospitalResolver(hospitalAAdapter domain.ExternalAPIAdapter) *HospitalResolver {
	return &HospitalResolver{
		hospitalAAdapter: hospitalAAdapter,
	}
}

func (r *HospitalResolver) Resolve(hospitalID string) (domain.ExternalAPIAdapter, error) {
	switch hospitalID {
	case "HOSP_A", "HN99999", "HN12345": // รองรับ Hospital ID ของ รพ. A
		return r.hospitalAAdapter, nil
	default:
		return nil, errors.New("ไม่พบรหัสโรงพยาบาลนี้ในระบบ หรือยังไม่รองรับครับ")
	}
}
