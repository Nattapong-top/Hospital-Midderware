package domain

type StaffRepository interface {
	FindByUsername(username string) (*Staff, error)
	Save(staff *Staff) error
}
