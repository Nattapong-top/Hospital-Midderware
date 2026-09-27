package infrastructure

import (
	"Hospital-Middleware/internal/domain"
)

type HospitalResolver struct {
	hospitalAAdapter domain.ExternalAPIAdapter
	hospitalBAdapter domain.ExternalAPIAdapter
}

func NewHospitalResolver(hospitalAAdapter domain.ExternalAPIAdapter, hospitalBAdapters ...domain.ExternalAPIAdapter) *HospitalResolver {
	var hospitalBAdapter domain.ExternalAPIAdapter
	if len(hospitalBAdapters) > 0 {
		hospitalBAdapter = hospitalBAdapters[0]
	} else {
		hospitalBAdapter = hospitalAAdapter
	}
	return &HospitalResolver{
		hospitalAAdapter: hospitalAAdapter,
		hospitalBAdapter: hospitalBAdapter,
	}
}

func (r *HospitalResolver) Resolve(hospitalID string) (domain.ExternalAPIAdapter, error) {
	switch hospitalID {
	case domain.HospitalHospA, domain.HospitalIDBangkok: // โรงพยาบาล A (Bangkok General Hospital)
		return r.hospitalAAdapter, nil
	case domain.HospitalHospB, domain.HospitalIDAgnos: // โรงพยาบาล B (Agnos Central Hospital)
		return r.hospitalBAdapter, nil
	default:
		return nil, domain.ErrHospitalNotFound
	}
}
