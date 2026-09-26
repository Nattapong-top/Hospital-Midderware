package infrastructure

import (
	"errors"

	"Hospital-Midderware/internal/domain"
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
	case "HOSP_A", "HN99999": // รองรับ Hospital ID ของ รพ. A
		return r.hospitalAAdapter, nil
	default:
		// ถ้ามี Hospital B เพิ่มในอนาคต สามารถมาสวิตช์เพิ่มตรงนี้ได้ง่ายๆ
		return nil, errors.New("unsupported or unknown hospital id")
	}
}
