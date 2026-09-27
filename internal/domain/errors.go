package domain

import "errors"

// Sentinel Errors สำหรับขอบเขต Domain และ Application
var (
	ErrUserAlreadyExists  = errors.New("username นี้มีอยู่ในระบบแล้ว")
	ErrInvalidCredentials = errors.New("ข้อมูลการเข้าสู่ระบบไม่ถูกต้อง")
	ErrInvalidCriteria    = errors.New("กรุณาระบุเงื่อนไขในการค้นหาอย่างน้อยหนึ่งรายการ")
	ErrHospitalNotFound   = errors.New("ไม่พบรหัสโรงพยาบาลนี้ในระบบ หรือยังไม่รองรับครับ")
)
