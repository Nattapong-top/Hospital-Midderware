package domain

// ExternalAPIAdapter เป็น Interface สำหรับการเชื่อมต่อยิงค้นหาข้อมูลผู้ป่วยกับ External HIS API
type ExternalAPIAdapter interface {
	Search(criteria SearchCriteria) (*PatientDTO, error)
}

// HospitalAPIResolver เป็น Interface สำหรับ Resolve หา ExternalAPIAdapter ตาม hospital_id
type HospitalAPIResolver interface {
	Resolve(hospitalID string) (ExternalAPIAdapter, error)
}
