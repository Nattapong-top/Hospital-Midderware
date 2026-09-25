package domain

import (
	"errors"
	"strings"
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
