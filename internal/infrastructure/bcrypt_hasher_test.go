package infrastructure

import "testing"

func TestBcryptHasher_HashAndCompare_Success(t *testing.T) {
	hasher := NewBcryptHasher()
	password := "khonnarak555"

	hashedPassword, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("คาดว่า Hash รหัสผ่านสำเร็จ แต่เกิด error: %v", err)
	}

	if hashedPassword == password {
		t.Fatalf("คาดว่ารหัสผ่านที่ Hash แล้วต้องไม่เหมือนเดิม")
	}

	isMatched := hasher.Compare(hashedPassword, password)
	if !isMatched {
		t.Error("คาดว่ารหัสผ่านต้องตรงกัน แต่ Compare คืนค่า false")
	}
}

func TestBcryptHasher_Compare_WrongPassword_ShouldFail(t *testing.T) {
	hasher := NewBcryptHasher()
	password := "khonnarak555"
	wrongPassword := "khonnarak888"

	hashedPassword, _ := hasher.Hash(password)

	isMatched := hasher.Compare(hashedPassword, wrongPassword)
	if isMatched {
		t.Error("คาดว่ารหัสผ่านผิด ต้องคืนค่า false แต่ได้ true")
	}
}
