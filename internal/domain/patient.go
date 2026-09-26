package domain

import "errors"

// SearchCriteria รับเงื่อนไขการค้นหาผู้ป่วย
type SearchCriteria struct {
	NationalID string
	PassportID string
	PatientHN  string
}

// Validate ตรวจสอบว่าต้องมีเงื่อนไขค้นหาอย่างน้อย 1 รายการ
func (s *SearchCriteria) Validate() error {
	if s.NationalID == "" && s.PassportID == "" && s.PatientHN == "" {
		return errors.New("at least one search criterion must be provided")
	}
	return nil
}

// PatientDTO แทนโครงสร้างข้อมูลผู้ป่วยที่ส่งมาจาก External HIS API
type PatientDTO struct {
	FirstNameTH  string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH   string `json:"last_name_th"`
	FirstNameEN  string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN   string `json:"last_name_en"`
	DateOfBirth  string `json:"date_of_birth"`
	PatientHN    string `json:"patient_hn"`
	NationalID   string `json:"national_id"`
	PassportID   string `json:"passport_id"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
}

// SearchPatientResult ผลลัพธ์จากการค้นหา
type SearchPatientResult struct {
	Patient *PatientDTO
}
