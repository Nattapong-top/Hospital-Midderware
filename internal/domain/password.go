package domain

import (
	"errors"
	"strings"
	"unicode"
)

type Password struct {
	value string
}

func NewPassword(value string) (Password, error) {

	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return Password{}, errors.New("กรุณากรอก password และห้ามกรอกเป็นค่าว่างครับ")
	}

	if len(trimmedValue) < 8 {
		return Password{}, errors.New("รหัสผ่านต้องมีความยาวอย่างน้อย 8 ตัวอักษร")
	}

	for _, ch := range trimmedValue {
		if ch > unicode.MaxASCII {
			return Password{}, errors.New("รหัสผ่านต้องใช้อักขระภาษาอังกฤษ ตัวเลข หรือสัญลักษณ์สากลเท่านั้น")
		}
	}

	return Password{
		value: value,
	}, nil
}

func (p Password) Value() string { return p.value }
