package domain

// PatientToDTO constructs a PatientDTO from raw string attributes.
func PatientToDTO(hn, nationalID, passportID, firstNameTH, lastNameTH, firstNameEN, lastNameEN, gender, dob, phone, email string) PatientDTO {
	return PatientDTO{
		PatientHN:   hn,
		NationalID:  nationalID,
		PassportID:  passportID,
		FirstNameTH: firstNameTH,
		LastNameTH:  lastNameTH,
		FirstNameEN: firstNameEN,
		LastNameEN:  lastNameEN,
		Gender:      gender,
		DateOfBirth: dob,
		PhoneNumber: phone,
		Email:       email,
	}
}

// StaffSummary represents a mapped representation of a staff member.
type StaffSummary struct {
	Username   string
	HospitalID string
}

// StaffToSummary maps a Staff domain entity to a StaffSummary.
func StaffToSummary(staff *Staff) StaffSummary {
	if staff == nil {
		return StaffSummary{}
	}
	return StaffSummary{
		Username:   staff.Username.Value(),
		HospitalID: staff.HospitalId.Value(),
	}
}
