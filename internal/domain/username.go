package domain

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type Username struct {
	value string
}

func NewUsername(value string) (Username, error) {
	trimmedValue := strings.TrimSpace(value)

	if trimmedValue == "" {
		return Username{}, errors.New("กรุณากรอก Username และห้ามเป็นค่าว่างครับ")
	}

	if utf8.RuneCountInString(trimmedValue) < 6 {
		return Username{}, errors.New("username ต้องมีความยาวอย่างน้อย 6 ตัวอักษร")
	}

	return Username{
		value: trimmedValue,
	}, nil
}

func (u Username) Value() string {
	return u.value
}
