package infrastructure

import (
	"Hospital-Midderware/internal/domain"
	"context"
	"database/sql"
	"errors"
)

type PostgresStaffRepository struct {
	db *sql.DB
}

func NewPostgresStaffRepository(db *sql.DB) *PostgresStaffRepository {
	return &PostgresStaffRepository{
		db: db,
	}
}

func (r *PostgresStaffRepository) Save(staff *domain.Staff) error {
	query := `
			INSERT INTO staffs (username, password, hospital_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (username) DO UPDATE
			SET password = EXCLUDED.password,
				hospital_id = EXCLUDED.hospital_id,
				updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(
		query,
		staff.Username.Value(),
		staff.Password.Value(),
		staff.HospitalId.Value(),
	)

	return err
}

func (r *PostgresStaffRepository) FindByUsername(username string) (*domain.Staff, error) {
	query := `
			SELECT username, password, hospital_id
			FROM staffs
			WHERE username = $1
	`

	var dbUsername, dbPassword, dbHospitalId string

	err := r.db.QueryRow(query, username).Scan(&dbUsername, &dbPassword, &dbHospitalId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ไม่พบ staff")
		}
		return nil, err
	}

	staff, err := domain.CreateStaff(dbUsername, dbPassword, dbHospitalId)
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (r *PostgresStaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	query := `
		INSERT INTO staffs (username, password, hospital_id)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		staff.Username.Value(),
		staff.Password.Value(),
		staff.HospitalId.Value(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresStaffRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM staffs WHERE username = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
