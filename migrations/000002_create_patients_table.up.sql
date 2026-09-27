-- 1. สร้างตาราง patients (ข้อมูลผู้ป่วย)
CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    patient_hn VARCHAR(50) NOT NULL UNIQUE,
    hospital_id VARCHAR(50) NOT NULL,
    national_id VARCHAR(20),
    passport_id VARCHAR(20),
    first_name_th VARCHAR(100),
    middle_name_th VARCHAR(100),
    last_name_th VARCHAR(100),
    first_name_en VARCHAR(100),
    middle_name_en VARCHAR(100),
    last_name_en VARCHAR(100),
    date_of_birth DATE,
    gender VARCHAR(10),
    phone_number VARCHAR(20),
    email VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_patients_hospital
    FOREIGN KEY (hospital_id)
    REFERENCES hospitals(id)
    ON DELETE RESTRICT
);

-- 2. สร้าง Indexes สำหรับเพิ่มความเร็วในการค้นหาผู้ป่วย
CREATE INDEX IF NOT EXISTS idx_patients_national_id ON patients(national_id);
CREATE INDEX IF NOT EXISTS idx_patients_passport_id ON patients(passport_id);
CREATE INDEX IF NOT EXISTS idx_patients_phone ON patients(phone_number);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_id ON patients(hospital_id);
