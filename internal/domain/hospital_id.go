package domain

import (
	"errors"
	"strings"
)

const (
	HospitalIDBangkok = "HN12345"
	HospitalIDAgnos   = "HN99999"
	HospitalHospA     = "HOSP_A"
	HospitalHospB     = "HOSP_B"
	HospitalHospC     = "HOSP_C"
)

type HospitalId struct {
	value string
}

func NewHospitalId(value string) (HospitalId, error) {

	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return HospitalId{}, errors.New("กรุณากรอก HospitalID ด้วยครับ")
	}

	return HospitalId{
		value: trimmedValue,
	}, nil
}

func (h HospitalId) Value() string { return h.value }
