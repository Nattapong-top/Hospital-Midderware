package infrastructure

import (
	"Hospital-Middleware/internal/domain"
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
		UPDATE staffs
		SET password = $2,
			hospital_id = $3,
			current_version = $4,
			previous_version = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE username = $1 AND current_version = $6
	`

	result, err := r.db.Exec(
		query,
		staff.Username.Value(),
		staff.Password.Value(),
		staff.HospitalId.Value(),
		staff.Version.CurrentNumber(),
		staff.Version.PreviousNumber(),
		staff.Version.PreviousNumber(),
	)
	if err != nil {
		return errors.New("เกิดข้อผิดพลาดในการอัปเดตข้อมูลพนักงานครับ: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("เกิดข้อผิดพลาดในการตรวจสอบผลการอัปเดตครับ: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("ข้อมูลพนักงานถูกแก้ไขโดยผู้อื่นแล้ว กรุณาลองใหม่อีกครั้งครับ")
	}

	return nil
}

func (r *PostgresStaffRepository) FindByUsername(username string) (*domain.Staff, error) {
	query := `
			SELECT username, password, hospital_id, current_version, previous_version
			FROM staffs
			WHERE username = $1
	`

	var dbUsername, dbPassword, dbHospitalId string
	var dbCurVer, dbPrevVer int

	err := r.db.QueryRow(query, username).Scan(&dbUsername, &dbPassword, &dbHospitalId, &dbCurVer, &dbPrevVer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ไม่พบข้อมูลพนักงานในระบบครับ")
		}
		return nil, errors.New("เกิดข้อผิดพลาดในการค้นหาข้อมูลพนักงานครับ: " + err.Error())
	}

	staff, err := domain.RestoreStaff(dbUsername, dbPassword, dbHospitalId, dbCurVer, dbPrevVer)
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (r *PostgresStaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	query := `
		INSERT INTO staffs (username, password, hospital_id, current_version, previous_version)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		staff.Username.Value(),
		staff.Password.Value(),
		staff.HospitalId.Value(),
		staff.Version.CurrentNumber(),
		staff.Version.PreviousNumber(),
	)
	if err != nil {
		return errors.New("เกิดข้อผิดพลาดในการสร้างบัญชีพนักงานครับ: " + err.Error())
	}
	return nil
}

func (r *PostgresStaffRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM staffs WHERE username = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, errors.New("เกิดข้อผิดพลาดในการตรวจสอบชื่อผู้ใช้งานครับ: " + err.Error())
	}

	return exists, nil
}
