package application

import (
	"Hospital-Midderware/internal/domain"
	"Hospital-Midderware/internal/infrastructure"
	"testing"
)

func TestAuthService_Login_Success(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider("secret-key-5555")

	hashedPassword, _ := hasher.Hash("RakTukKhon5555")
	staff, _ := domain.CreateStaff("nattapong", hashedPassword, "HN99999")
	err := repo.Save(staff)
	if err != nil {
		t.Fatalf("บันทึก staff ลง repository ไม่สำเร็จ: %v", err)
	}

	authService := NewAuthService(repo, hasher, tokenProvider)

	token, err := authService.Login("nattapong", "RakTukKhon5555")

	if err != nil {
		t.Fatalf("คาดว่า login สำเร็จ แต่ได้ error: %v", err)
	}

	if token == "" {
		t.Fatalf("คาดว่าต้องได้ Token แต่ได้ข้อความว่างเปล่า")
	}
}

func TestAuthService_Login_UserNotFound_ShouldFail(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider("secret-key-5555")

	authService := NewAuthService(repo, hasher, tokenProvider)

	_, err := authService.Login("unknown_user", "RakTukKhon5555")

	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อไม่พบผู้ใช้งาน แต่กลับผ่าน")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("คาดหวัง error %q แต่ได้ %q", "invalid credentials", err.Error())
	}
}

func TestAuthService_Login_WrongPassword_ShouldFail(t *testing.T) {
	repo := infrastructure.NewInMemoryStaffRepository()
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider("secret-key-5555")

	hashedPassword, _ := hasher.Hash("RakTukKhon5555")
	staff, _ := domain.CreateStaff("nattapong", hashedPassword, "HN99999")
	err := repo.Save(staff)
	if err != nil {
		t.Fatalf("บันทึก staff ลง repository ไม่สำเร็จ: %v", err)
	}

	authService := NewAuthService(repo, hasher, tokenProvider)

	_, err = authService.Login("nattapong", "wrong_password")

	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อพิมพ์รหัสผ่านผิด แต่กลับผ่าน")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("คาดหวัง error %q แต่ได้ %q", "invalid credentials", err.Error())
	}
}

func TestAuthService_Login_MissingHospitalId_ShouldFail(t *testing.T) {
	tokenProvider := infrastructure.NewJWTTokenProvider("secret-key-5555")

	_, err := tokenProvider.GenerateToken(nil)
	if err == nil {
		t.Fatal("คาดว่าจะเกิด error เมื่อ staff เป็น nil แต่กลับผ่าน")
	}
}
