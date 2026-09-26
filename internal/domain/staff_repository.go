package domain

import "context"

type StaffRepository interface {
	FindByUsername(username string) (*Staff, error)
	Save(staff *Staff) error
	Create(ctx context.Context, staff *Staff) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
