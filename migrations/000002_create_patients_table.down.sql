-- Drop indexes and table for patients
DROP INDEX IF EXISTS idx_patients_hospital_id;
DROP INDEX IF EXISTS idx_patients_phone;
DROP INDEX IF EXISTS idx_patients_passport_id;
DROP INDEX IF EXISTS idx_patients_national_id;
DROP TABLE IF EXISTS patients;
