package domain

import "testing"

func TestNewPassword_Success(t *testing.T) {
	password := "password"
	passwordVo, err := NewPassword(password)

	if err != nil {
		t.Fatalf("คาดว่าไมมี error แต่มี err: %v", err)
	}

	if password != passwordVo.Value() {
		t.Errorf(
			"คาดหวังให้ password %q, แต่ได้รับ %q",
			password,
			passwordVo.Value(),
		)
	}
}

func TestNewPassword_entry(t *testing.T) {
	password := "        "
	_, err := NewPassword(password)
	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อ password ว่าง")
	}
}

func TestNewPassword_MinLength_LessThan8Chars(t *testing.T) {
	password := "1234567"
	_, err := NewPassword(password)
	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมือความยาว username น้อยกว่า 8 ตัวอักษร")
	}
}

func TestNewPassword_ThaiLanguage_ShouldFail(t *testing.T) {
	password := "testทดสอบ123"
	_, err := NewPassword(password)
	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อใช้ภาษาไทยสร้าง password")
	}
}
