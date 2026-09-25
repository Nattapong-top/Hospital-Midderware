package domain

type Staff struct {
	Username   Username
	Password   Password
	HospitalId HospitalId
}

func CreateStaff(username, password, hospitalId string) (*Staff, error) {
	usernameVo, err := NewUsername(username)
	if err != nil {
		return nil, err
	}

	passwordVo, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	hospitalIdVo, err := NewHospitalId(hospitalId)
	if err != nil {
		return nil, err
	}

	return &Staff{
		Username:   usernameVo,
		Password:   passwordVo,
		HospitalId: hospitalIdVo,
	}, nil
}
