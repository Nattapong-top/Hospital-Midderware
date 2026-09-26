-- 1. สร้างตารางโรงพยาบาล (Master Data)
CREATE TABLE IF NOT EXISTS hospitals (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. สร้างตารางพนักงาน (ผูก FK หา hospitals)
CREATE TABLE IF NOT EXISTS staffs (
    id SERIAL PRIMARY KEY,
    username VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    hospital_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

-- Foreign Key: ห้ามใส่ hospital_id ที่ไม่มีอยู่ในตาราง hospitals
    CONSTRAINT fk_staffs_hospital
    FOREIGN KEY (hospital_id)
    REFERENCES hospitals(id)
    ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_staffs_username ON staffs(username);

-- 3. Seed ข้อมูลโรงพยาบาลเริ่มต้นสำหรับทดสอบ (Dummy Data)
INSERT INTO hospitals (id, name) VALUES
    ('HN99999', 'Agnos Central Hospital'),
    ('HN12345', 'Bangkok General Hospital')
ON CONFLICT (id) DO NOTHING;
