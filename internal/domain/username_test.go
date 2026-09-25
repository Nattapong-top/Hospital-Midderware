package domain

import "testing"

func TestNewUsername_Success(t *testing.T) {
	username := "nattap"
	usernameVo, err := NewUsername(username)

	if err != nil {
		t.Fatalf("คาดว่าไม่มี error แต่มี err: %v", err)
	}

	if usernameVo.Value() != username {
		t.Errorf(
			"คาดหวังให้ username %q, ได้รับ %q",
			username,
			usernameVo.Value(),
		)
	}
}

func TestNewUsername_entry(t *testing.T) {
	username := ""
	usernameVo, err := NewUsername(username)

	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อ username ว่าง")
	}

	if usernameVo.Value() != "" {
		t.Errorf("คาดว่า username จะว่าง แต่ได้ %q", usernameVo.Value())
	}
}

func TestUsername_MinLen_LessThan6Chars(t *testing.T) {
	username := "natta"
	_, err := NewUsername(username)
	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อความยาว username น้อยกว่า 6 ตัวอักษร")
	}
}
